package usecase

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/m11ano/e"
	productscl "github.com/m11ano/mipt-webdev-course/backend/clients/clgrpc/pkg/products"
	productsgcl "github.com/m11ano/mipt-webdev-course/backend/services/orders/internal/clients/grpc/products"
	"github.com/m11ano/mipt-webdev-course/backend/services/orders/internal/domain"
	"github.com/m11ano/mipt-webdev-course/backend/services/orders/internal/infra/config"
	"github.com/m11ano/mipt-webdev-course/backend/services/orders/internal/usecase/uctypes"
	productstc "github.com/m11ano/mipt-webdev-course/backend/temporal-app/pkg/workers/products/client"
	"github.com/samber/lo"
	"github.com/shopspring/decimal"
)

var ErrOrderInvalidProducts = e.NewErrorFrom(e.ErrBadRequest).SetMessage("invalid products")
var ErrOrderInvalidProductsQuantity = e.NewErrorFrom(e.ErrBadRequest).SetMessage("invalid products quantity")

type OrderPartUpdateData struct {
	ClientName      *string
	ClientSurname   *string
	ClientEmail     *string
	ClientPhone     *string
	DeliveryAddress *string
}

type OrderListSortField int

const (
	OrderListSortFieldID OrderListSortField = iota
)

type OrderListSort struct {
	Field  OrderListSortField
	IsDesc bool
}

type OrderListOptions struct {
	IDs         *[]int64
	OnlyCreated *bool
	Sort        *[]OrderListSort
}

type OrderCreateIn struct {
	Details  OrderDataDetailsIn
	Products []OrderProductIn
}

type OrderUpdateIn struct {
	Details  OrderDataDetailsIn
	Products []OrderProductWithPrice
}

type OrderDataDetailsIn struct {
	ClientName      string
	ClientSurname   string
	ClientEmail     string
	ClientPhone     string
	DeliveryAddress string
}

type OrderProductIn struct {
	ID       int64
	Quantity int32
}

type OrderProductWithPrice struct {
	ID       int64
	Quantity int32
	Price    decimal.Decimal
}

type OrderOneFullOut struct {
	Order    *domain.Order
	Products []OrderProductWithPrice
}

type SetOrderCompositionIn struct {
	OrderID  int64
	Products *[]OrderProductWithPrice
	Status   *domain.OrderStatus
}

//go:generate mockery --name=Order --output=../../tests/mocks --case=underscore
type Order interface {
	FindPagedList(ctx context.Context, listOptions OrderListOptions, queryParams *uctypes.QueryGetListParams) (out []*domain.Order, total int64, err error)
	FindList(ctx context.Context, listOptions OrderListOptions, queryParams *uctypes.QueryGetListParams) (out []*domain.Order, err error)
	FindOneFullByID(ctx context.Context, id int64, queryParams *uctypes.QueryGetOneParams) (out *OrderOneFullOut, err error)
	Create(ctx context.Context, input OrderCreateIn) (order *domain.Order, err error)
	Update(ctx context.Context, orderID int64, input OrderUpdateIn) (err error)
	SetOrderComposition(ctx context.Context, input SetOrderCompositionIn) (err error)
	RemoveOrderIfNew(ctx context.Context, orderID int64) (err error)
	SetStatus(ctx context.Context, orderID int64, status domain.OrderStatus) (err error)
}

//go:generate mockery --name=OrderRepository --output=../../tests/mocks --case=underscore
type OrderRepository interface {
	FindList(ctx context.Context, listOptions OrderListOptions, queryParams *uctypes.QueryGetListParams) (items []*domain.Order, err error)
	FindPagedList(ctx context.Context, listOptions OrderListOptions, queryParams *uctypes.QueryGetListParams) (items []*domain.Order, total int64, err error)
	FindOneByID(ctx context.Context, id int64, queryParams *uctypes.QueryGetOneParams) (order *domain.Order, err error)
	Create(ctx context.Context, item *domain.Order) (err error)
	Update(ctx context.Context, item *domain.Order) (err error)
	PartUpdateByList(ctx context.Context, updateData OrderPartUpdateData, listOptions OrderListOptions, withDeleted bool) (err error)
	PartUpdateByID(ctx context.Context, updateData OrderPartUpdateData, id int64) (err error)
	DeleteByList(ctx context.Context, listOptions OrderListOptions, isHardRemove bool) (err error)
}

type OrderInpl struct {
	logger         *slog.Logger
	config         config.Config
	repo           OrderRepository
	txManager      *manager.Manager
	productsGCl    *productsgcl.ClientConn
	productsTCl    productstc.Client
	orderProductUC OrderProduct
}

func NewOrderInpl(logger *slog.Logger, config config.Config, txManager *manager.Manager, repo OrderRepository, productsGCl *productsgcl.ClientConn, productsTCl productstc.Client, orderProductUC OrderProduct) *OrderInpl {
	uc := &OrderInpl{
		logger:         logger,
		config:         config,
		txManager:      txManager,
		repo:           repo,
		productsGCl:    productsGCl,
		productsTCl:    productsTCl,
		orderProductUC: orderProductUC,
	}
	return uc
}

func (uc *OrderInpl) FindPagedList(ctx context.Context, listOptions OrderListOptions, queryParams *uctypes.QueryGetListParams) ([]*domain.Order, int64, error) {

	list, total, err := uc.repo.FindPagedList(ctx, listOptions, queryParams)
	if err != nil {
		return nil, 0, err
	}

	return list, total, nil
}

func (uc *OrderInpl) FindList(ctx context.Context, listOptions OrderListOptions, queryParams *uctypes.QueryGetListParams) ([]*domain.Order, error) {

	list, err := uc.repo.FindList(ctx, listOptions, queryParams)
	if err != nil {
		return nil, err
	}

	return list, nil
}

func (uc *OrderInpl) FindOneFullByID(ctx context.Context, id int64, queryParams *uctypes.QueryGetOneParams) (*OrderOneFullOut, error) {
	order, err := uc.repo.FindOneByID(ctx, id, queryParams)
	if err != nil {
		return nil, err
	}

	products, err := uc.orderProductUC.FindList(ctx, OrderProductListOptions{
		OrderID: lo.ToPtr(id),
	}, nil)
	if err != nil {
		return nil, err
	}

	out := &OrderOneFullOut{
		Order:    order,
		Products: make([]OrderProductWithPrice, len(products)),
	}

	for i, product := range products {
		out.Products[i] = OrderProductWithPrice{
			ID:       product.ProductID,
			Quantity: product.Quantity,
			Price:    product.Price,
		}
	}

	return out, nil
}

func (uc *OrderInpl) Create(ctx context.Context, input OrderCreateIn) (*domain.Order, error) {

	productIDs := make([]int64, len(input.Products))
	for i, item := range input.Products {
		productIDs[i] = item.ID
	}
	productIDs = lo.Uniq(productIDs)

	if len(productIDs) == 0 || len(productIDs) != len(input.Products) {
		return nil, ErrOrderInvalidProducts
	}

	products, err := uc.productsGCl.Client.GetProductsByIds(ctx, productIDs)
	if err != nil {
		return nil, err
	}

	if len(products) != len(productIDs) {
		return nil, ErrOrderInvalidProducts
	}

	err = uc.checkStockAvailable(products, input.Products)
	if err != nil {
		return nil, err
	}

	order := domain.NewOrder(0)
	order.ClientName = input.Details.ClientName
	order.ClientSurname = input.Details.ClientSurname
	order.ClientEmail = input.Details.ClientEmail
	order.ClientPhone = input.Details.ClientPhone
	order.DeliveryAddress = input.Details.DeliveryAddress

	err = uc.repo.Create(ctx, order)
	if err != nil {
		return nil, err
	}

	orderSum := decimal.Zero
	for _, product := range input.Products {
		productItem, ok := lo.Find(products, func(item *productscl.ProductListItem) bool {
			return item.ID == product.ID
		})
		if !ok {
			return nil, e.ErrInternal
		}

		orderProduct, err := domain.NewOrderProduct(order.ID, product.ID, product.Quantity, productItem.Price)
		if err != nil {
			return nil, err
		}

		err = uc.orderProductUC.Create(ctx, orderProduct)
		if err != nil {
			return nil, err
		}

		orderSum = orderSum.Add(productItem.Price.Mul(decimal.NewFromInt(int64(product.Quantity))))
	}

	err = order.SetOrderSum(orderSum)
	if err != nil {
		return nil, err
	}

	err = uc.repo.Update(ctx, order)
	if err != nil {
		return nil, err
	}

	//Запускаем воркфлоу
	ctxWithTimeout, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	ordersList := make([]productstc.OrderProductsItem, len(input.Products))
	for i, item := range input.Products {
		productItem, ok := lo.Find(products, func(product *productscl.ProductListItem) bool {
			return product.ID == item.ID
		})
		if !ok {
			return nil, e.ErrInternal
		}

		ordersList[i] = productstc.OrderProductsItem{
			ProductID: item.ID,
			Quantity:  item.Quantity,
			Price:     productItem.Price,
		}
	}

	flowIn := productstc.SetOrderProductsAndStatusIn{
		OrderID:       order.ID,
		OrderProducts: &ordersList,
		OrderStatus:   lo.ToPtr(domain.OrderStatusCreated.String()),
	}

	err = uc.productsTCl.SetOrderProductsAndStatus(ctxWithTimeout, flowIn)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return nil, e.NewErrorFrom(e.ErrServiceUnavailable)
		}

		if errors.Is(err, productstc.ErrSetOrderProductsAndStatusCantReserve) {
			return nil, ErrOrderInvalidProductsQuantity
		}

		return nil, err
	}

	return order, nil
}

func (uc *OrderInpl) Update(ctx context.Context, orderID int64, input OrderUpdateIn) error {

	productIDs := make([]int64, len(input.Products))
	for i, item := range input.Products {
		productIDs[i] = item.ID
	}
	productIDs = lo.Uniq(productIDs)

	if len(productIDs) == 0 || len(productIDs) != len(input.Products) {
		return ErrOrderInvalidProducts
	}

	products, err := uc.productsGCl.Client.GetProductsByIds(ctx, productIDs)
	if err != nil {
		return err
	}

	if len(products) != len(productIDs) {
		return ErrOrderInvalidProducts
	}

	for _, product := range products {
		if !product.IsPublished {
			return ErrOrderInvalidProducts
		}
	}

	order, err := uc.repo.FindOneByID(ctx, orderID, nil)
	if err != nil {
		return err
	}

	if order.Status == domain.OrderStatusCanceled {
		return e.NewErrorFrom(e.ErrBadRequest).SetMessage("order is canceled")
	}

	if order.Status == domain.OrderStatusFinished {
		return e.NewErrorFrom(e.ErrBadRequest).SetMessage("order is finished")
	}

	//Запускаем воркфлоу
	ctxWithTimeout, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	ordersList := make([]productstc.OrderProductsItem, len(input.Products))
	for i, item := range input.Products {
		ordersList[i] = productstc.OrderProductsItem{
			ProductID: item.ID,
			Quantity:  item.Quantity,
			Price:     item.Price,
		}
	}

	flowIn := productstc.SetOrderProductsAndStatusIn{
		OrderID:       order.ID,
		OrderProducts: &ordersList,
	}

	err = uc.productsTCl.SetOrderProductsAndStatus(ctxWithTimeout, flowIn)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return e.NewErrorFrom(e.ErrServiceUnavailable)
		}

		if errors.Is(err, productstc.ErrSetOrderProductsAndStatusCantReserve) {
			return ErrOrderInvalidProductsQuantity
		}

		return err
	}

	err = uc.repo.PartUpdateByID(ctx, OrderPartUpdateData{
		ClientName:      lo.ToPtr(input.Details.ClientName),
		ClientSurname:   lo.ToPtr(input.Details.ClientSurname),
		ClientEmail:     lo.ToPtr(input.Details.ClientEmail),
		ClientPhone:     lo.ToPtr(input.Details.ClientPhone),
		DeliveryAddress: lo.ToPtr(input.Details.DeliveryAddress),
	}, order.ID)
	if err != nil {
		return err
	}

	return nil
}

func (uc *OrderInpl) checkStockAvailable(products []*productscl.ProductListItem, orderProducts []OrderProductIn) error {
	details := []string{}

	for _, orderProduct := range orderProducts {
		if orderProduct.Quantity < 1 {
			return e.NewErrorFrom(e.ErrBadRequest).SetMessage(fmt.Sprintf("invalid quantity in #%d product", orderProduct.ID))
		}

		product, isFound := lo.Find(products, func(item *productscl.ProductListItem) bool {
			return item.ID == orderProduct.ID
		})

		if !isFound || !product.IsPublished {
			return ErrOrderInvalidProducts
		}

		if product.StockAvailable < orderProduct.Quantity {
			details = append(details, fmt.Sprintf(`Доступный для заказа остаток по товару "%s": %d шт.`, product.Name, product.StockAvailable))
		}
	}

	if len(details) > 0 {
		return e.NewErrorFrom(ErrOrderInvalidProductsQuantity).AddDetails(details)
	}

	return nil
}

func (uc *OrderInpl) SetOrderComposition(ctx context.Context, input SetOrderCompositionIn) error {

	if input.Products == nil && input.Status == nil {
		return e.NewErrorFrom(e.ErrBadRequest).SetMessage("products or status must be set")
	}

	err := uc.txManager.Do(ctx, func(ctx context.Context) error {
		order, err := uc.repo.FindOneByID(ctx, input.OrderID, &uctypes.QueryGetOneParams{
			ForUpdate: true,
		})
		if err != nil {
			return err
		}

		if order.Status == domain.OrderStatusCanceled || order.Status == domain.OrderStatusFinished {
			return e.NewErrorFrom(e.ErrBadRequest).SetMessage("order is canceled or finished")
		}

		skeepProducts := false

		if input.Status != nil {
			err = order.SetStatus(*input.Status)
			if err != nil {
				return err
			}

			if order.Status == domain.OrderStatusCanceled {
				skeepProducts = true
			}
		}

		if input.Products != nil && !skeepProducts {

			err = uc.orderProductUC.DeleteByOrderID(ctx, input.OrderID)
			if err != nil {
				return err
			}

			orderSum := decimal.Zero
			for _, item := range *input.Products {
				orderProduct, err := domain.NewOrderProduct(input.OrderID, item.ID, item.Quantity, item.Price)
				if err != nil {
					return err
				}

				err = uc.orderProductUC.Create(ctx, orderProduct)
				if err != nil {
					return err
				}

				orderSum = orderSum.Add(item.Price.Mul(decimal.NewFromInt(int64(item.Quantity))))
			}

			err = order.SetOrderSum(orderSum)
			if err != nil {
				return err
			}
		}

		err = uc.repo.Update(ctx, order)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (uc *OrderInpl) RemoveOrderIfNew(ctx context.Context, orderID int64) error {

	err := uc.txManager.Do(ctx, func(ctx context.Context) error {
		order, err := uc.repo.FindOneByID(ctx, orderID, &uctypes.QueryGetOneParams{
			ForUpdate: true,
		})
		if err != nil {
			return err
		}

		if order.Status != domain.OrderStatusNew {
			return e.NewErrorFrom(e.ErrBadRequest).SetMessage("order is not new")
		}

		err = uc.orderProductUC.DeleteByOrderID(ctx, orderID)
		if err != nil {
			return err
		}

		err = uc.repo.DeleteByList(ctx, OrderListOptions{
			IDs: lo.ToPtr([]int64{orderID}),
		}, true)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (uc *OrderInpl) SetStatus(ctx context.Context, orderID int64, status domain.OrderStatus) error {

	if status == domain.OrderStatusNew {
		return e.NewErrorFrom(e.ErrBadRequest).SetMessage("cant set status")
	}

	order, err := uc.repo.FindOneByID(ctx, orderID, &uctypes.QueryGetOneParams{
		ForUpdate: true,
	})
	if err != nil {
		return err
	}

	if order.Status == status {
		return nil
	}

	err = order.SetStatus(status)
	if err != nil {
		return err
	}

	//Запускаем воркфлоу и не ждем результат
	flowIn := productstc.SetOrderProductsAndStatusIn{
		NotWait:     true,
		OrderID:     order.ID,
		OrderStatus: lo.ToPtr(order.Status.String()),
	}

	//Если заказ отменен, то отменяем списание товаров
	if order.Status == domain.OrderStatusCanceled {
		flowIn.OrderProducts = lo.ToPtr([]productstc.OrderProductsItem{})
	}

	err = uc.productsTCl.SetOrderProductsAndStatus(ctx, flowIn)
	if err != nil {
		return err
	}

	return nil
}
