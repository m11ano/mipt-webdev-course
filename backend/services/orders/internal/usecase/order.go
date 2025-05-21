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
	IDs  *[]int64
	Sort *[]OrderListSort
}

type OrderCreateIn struct {
	Details  OrderCreateInDetails
	Products []OrderProductIn
}

type OrderCreateInDetails struct {
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

type OrderProductWithPriceIn struct {
	ID       int64
	Quantity int32
	Price    decimal.Decimal
}

type OrderOneFullOut struct {
	Order *domain.Order
}

type OrderFullOut struct {
	Order *domain.Order
}

type SetOrderCompositionIn struct {
	OrderID  int64
	Products []OrderProductWithPriceIn
}

//go:generate mockery --name=Order --output=../../tests/mocks --case=underscore
type Order interface {
	FindFullPagedList(ctx context.Context, listOptions OrderListOptions, queryParams *uctypes.QueryGetListParams) (out []*OrderFullOut, total int64, err error)
	FindFullList(ctx context.Context, listOptions OrderListOptions, queryParams *uctypes.QueryGetListParams) (out []*OrderFullOut, err error)
	FindOneFullByID(ctx context.Context, id int64, queryParams *uctypes.QueryGetOneParams) (out *OrderOneFullOut, err error)
	Create(ctx context.Context, input OrderCreateIn) (order *domain.Order, err error)
	SetOrderComposition(ctx context.Context, input SetOrderCompositionIn) (err error)
	RemoveNewOrder(ctx context.Context, orderID int64) (err error)
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

func (uc *OrderInpl) FindFullPagedList(ctx context.Context, listOptions OrderListOptions, queryParams *uctypes.QueryGetListParams) ([]*OrderFullOut, int64, error) {

	list, total, err := uc.repo.FindPagedList(ctx, listOptions, queryParams)
	if err != nil {
		return nil, 0, err
	}

	result := make([]*OrderFullOut, len(list))
	for i, item := range list {
		result[i] = &OrderFullOut{
			Order: item,
		}
	}

	return result, total, nil
}

func (uc *OrderInpl) FindFullList(ctx context.Context, listOptions OrderListOptions, queryParams *uctypes.QueryGetListParams) ([]*OrderFullOut, error) {

	list, err := uc.repo.FindList(ctx, listOptions, queryParams)
	if err != nil {
		return nil, err
	}

	result := make([]*OrderFullOut, len(list))
	for i, item := range list {
		result[i] = &OrderFullOut{
			Order: item,
		}
	}

	return result, nil
}

func (uc *OrderInpl) FindOneFullByID(ctx context.Context, id int64, queryParams *uctypes.QueryGetOneParams) (*OrderOneFullOut, error) {
	order, err := uc.repo.FindOneByID(ctx, id, queryParams)
	if err != nil {
		return nil, err
	}

	out := &OrderOneFullOut{
		Order: order,
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
	ctxWithTimeout, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	flowIn := productstc.SetOrderProductsAndBlockIn{
		OrderID:       order.ID,
		OrderProducts: make([]productstc.OrderProductsItem, len(input.Products)),
	}

	for i, item := range input.Products {
		productItem, ok := lo.Find(products, func(product *productscl.ProductListItem) bool {
			return product.ID == item.ID
		})
		if !ok {
			return nil, e.ErrInternal
		}

		flowIn.OrderProducts[i] = productstc.OrderProductsItem{
			ProductID: item.ID,
			Quantity:  item.Quantity,
			Price:     productItem.Price,
		}
	}

	err = uc.productsTCl.SetOrderProductsAndBlock(ctxWithTimeout, flowIn)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return nil, e.NewErrorFrom(e.ErrServiceUnavailable)
		}

		if errors.Is(err, productstc.ErrSetOrderProductsAndBlockCantReserve) {
			return nil, ErrOrderInvalidProductsQuantity
		}

		return nil, err
	}

	return order, nil
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

		if !isFound {
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

		if order.Status == domain.OrderStatusNew {
			order.Status = domain.OrderStatusCreated
		}

		err = uc.orderProductUC.DeleteByOrderID(ctx, input.OrderID)
		if err != nil {
			return err
		}

		orderSum := decimal.Zero
		for _, item := range input.Products {
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

func (uc *OrderInpl) RemoveNewOrder(ctx context.Context, orderID int64) error {

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
