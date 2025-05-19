package usecase

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/m11ano/e"
	"github.com/m11ano/mipt-webdev-course/backend/services/orders/internal/domain"
	"github.com/m11ano/mipt-webdev-course/backend/services/orders/internal/infra/config"
	"github.com/m11ano/mipt-webdev-course/backend/services/orders/internal/usecase/uctypes"
	"github.com/samber/lo"
)

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
	Products []OrderCreateInProduct
}

type OrderCreateInDetails struct {
	ClientName      string
	ClientSurname   string
	ClientEmail     string
	ClientPhone     string
	DeliveryAddress string
}

type OrderCreateInProduct struct {
	ID       int64
	Quantity int32
}

type OrderOneFullOut struct {
	Order *domain.Order
}

type OrderFullOut struct {
	Order *domain.Order
}

//go:generate mockery --name=Order --output=../../tests/mocks --case=underscore
type Order interface {
	FindFullPagedList(ctx context.Context, listOptions OrderListOptions, queryParams *uctypes.QueryGetListParams) (out []*OrderFullOut, total int64, err error)
	FindFullList(ctx context.Context, listOptions OrderListOptions, queryParams *uctypes.QueryGetListParams) (out []*OrderFullOut, err error)
	FindOneFullByID(ctx context.Context, id int64, queryParams *uctypes.QueryGetOneParams) (out *OrderOneFullOut, err error)
	Create(ctx context.Context, input OrderCreateIn) (order *domain.Order, err error)
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
	DeleteByList(ctx context.Context, listOptions OrderListOptions) (err error)
}

type OrderInpl struct {
	logger    *slog.Logger
	config    config.Config
	repo      OrderRepository
	txManager *manager.Manager
}

func NewOrderInpl(logger *slog.Logger, config config.Config, txManager *manager.Manager, repo OrderRepository) *OrderInpl {
	uc := &OrderInpl{
		logger:    logger,
		config:    config,
		txManager: txManager,
		repo:      repo,
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
		return nil, e.NewErrorFrom(e.ErrBadRequest).SetMessage("invalid products")
	}

	// Сразу предварительно до создания воркфлоу проверим наличие, в случае если товара нет - не будем создавать заведомо провальный воркфлоу
	fmt.Println(productIDs)

	order := domain.NewOrder(0)
	order.ClientName = input.Details.ClientName
	order.ClientSurname = input.Details.ClientSurname
	order.ClientEmail = input.Details.ClientEmail
	order.ClientPhone = input.Details.ClientPhone
	order.DeliveryAddress = input.Details.DeliveryAddress

	return order, nil
}
