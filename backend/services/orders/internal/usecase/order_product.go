package usecase

import (
	"context"
	"log/slog"

	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/m11ano/mipt-webdev-course/backend/services/orders/internal/domain"
	"github.com/m11ano/mipt-webdev-course/backend/services/orders/internal/infra/config"
	"github.com/m11ano/mipt-webdev-course/backend/services/orders/internal/usecase/uctypes"
)

type OrderProductListOptions struct {
	ProductID *int64
	OrderID   *int64
}

//go:generate mockery --name=OrderProduct --output=../../tests/mocks --case=underscore
type OrderProduct interface {
	FindList(ctx context.Context, listOptions OrderProductListOptions, queryParams *uctypes.QueryGetListParams) (items []*domain.OrderProduct, err error)
	Create(ctx context.Context, item *domain.OrderProduct) (err error)
	DeleteByOrderID(ctx context.Context, orderID int64) (err error)
}

//go:generate mockery --name=OrderProductRepository --output=../../tests/mocks --case=underscore
type OrderProductRepository interface {
	FindList(ctx context.Context, listOptions OrderProductListOptions, queryParams *uctypes.QueryGetListParams) (items []*domain.OrderProduct, err error)
	Create(ctx context.Context, item *domain.OrderProduct) (err error)
	DeleteByList(ctx context.Context, listOptions OrderProductListOptions) (err error)
}

type OrderProductRepositoryInpl struct {
	logger    *slog.Logger
	config    config.Config
	repo      OrderProductRepository
	txManager *manager.Manager
}

func NewOrderProductInpl(logger *slog.Logger, config config.Config, txManager *manager.Manager, repo OrderProductRepository) *OrderProductRepositoryInpl {
	uc := &OrderProductRepositoryInpl{
		logger:    logger,
		config:    config,
		txManager: txManager,
		repo:      repo,
	}
	return uc
}

func (uc *OrderProductRepositoryInpl) FindList(ctx context.Context, listOptions OrderProductListOptions, queryParams *uctypes.QueryGetListParams) ([]*domain.OrderProduct, error) {
	return uc.repo.FindList(ctx, listOptions, queryParams)
}

func (uc *OrderProductRepositoryInpl) Create(ctx context.Context, item *domain.OrderProduct) error {
	return uc.repo.Create(ctx, item)
}

func (uc *OrderProductRepositoryInpl) DeleteByOrderID(ctx context.Context, orderID int64) error {
	return uc.repo.DeleteByList(ctx, OrderProductListOptions{
		OrderID: &orderID,
	})
}
