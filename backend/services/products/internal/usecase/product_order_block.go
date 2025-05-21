package usecase

import (
	"context"
	"log/slog"

	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/m11ano/mipt-webdev-course/backend/services/products/internal/domain"
	"github.com/m11ano/mipt-webdev-course/backend/services/products/internal/infra/config"
	"github.com/m11ano/mipt-webdev-course/backend/services/products/internal/usecase/uctypes"
)

type ProductOrderBlockListOptions struct {
	ProductID *int64
	OrderID   *int64
}

type ProductOrderBlockComposition struct {
	ProductID int64
	Quantity  int32
}

//go:generate mockery --name=ProductOrderBlock --output=../../tests/mocks --case=underscore
type ProductOrderBlock interface {
	FindList(ctx context.Context, listOptions ProductOrderBlockListOptions, queryParams *uctypes.QueryGetListParams) (out []*domain.ProductOrderBlock, err error)
	Create(ctx context.Context, item *domain.ProductOrderBlock) (err error)
	ClearBlocksForOrder(ctx context.Context, orderID int64) (err error)
	GetOrderBlockedProducts(ctx context.Context, orderID int64) (items []*domain.ProductOrderBlock, err error)
	CheckBlockForProduct(ctx context.Context, productID int64) (result bool, err error)
}

//go:generate mockery --name=ProductOrderBlockRepository --output=../../tests/mocks --case=underscore
type ProductOrderBlockRepository interface {
	FindList(ctx context.Context, listOptions ProductOrderBlockListOptions, queryParams *uctypes.QueryGetListParams) (items []*domain.ProductOrderBlock, err error)
	Create(ctx context.Context, item *domain.ProductOrderBlock) (err error)
	DeleteByList(ctx context.Context, listOptions ProductOrderBlockListOptions) (err error)
}

type ProductOrderBlockInpl struct {
	logger    *slog.Logger
	config    config.Config
	repo      ProductOrderBlockRepository
	txManager *manager.Manager
}

func NewProductOrderBlockInpl(logger *slog.Logger, config config.Config, txManager *manager.Manager, repo ProductOrderBlockRepository) *ProductOrderBlockInpl {
	uc := &ProductOrderBlockInpl{
		logger:    logger,
		config:    config,
		txManager: txManager,
		repo:      repo,
	}
	return uc
}

func (uc *ProductOrderBlockInpl) FindList(ctx context.Context, listOptions ProductOrderBlockListOptions, queryParams *uctypes.QueryGetListParams) ([]*domain.ProductOrderBlock, error) {
	return uc.repo.FindList(ctx, listOptions, queryParams)
}

func (uc *ProductOrderBlockInpl) Create(ctx context.Context, item *domain.ProductOrderBlock) error {
	return uc.repo.Create(ctx, item)
}

func (uc *ProductOrderBlockInpl) ClearBlocksForOrder(ctx context.Context, orderID int64) error {
	return uc.repo.DeleteByList(ctx, ProductOrderBlockListOptions{
		OrderID: &orderID,
	})
}

func (uc *ProductOrderBlockInpl) GetOrderBlockedProducts(ctx context.Context, orderID int64) ([]*domain.ProductOrderBlock, error) {
	return uc.repo.FindList(ctx, ProductOrderBlockListOptions{
		OrderID: &orderID,
	}, &uctypes.QueryGetListParams{})
}

func (uc *ProductOrderBlockInpl) CheckBlockForProduct(ctx context.Context, productID int64) (bool, error) {
	check, err := uc.repo.FindList(ctx, ProductOrderBlockListOptions{
		ProductID: &productID,
	}, &uctypes.QueryGetListParams{
		Limit: 1,
	})
	if err != nil {
		return false, err
	}

	if len(check) > 0 {
		return true, nil
	}

	return false, nil
}
