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

//go:generate mockery --name=ProductOrderBlock --output=../../tests/mocks --case=underscore
type ProductOrderBlock interface {
	CheckBlockForProduct(ctx context.Context, productID int64) (result bool, err error)
	CreateBlockForProduct(ctx context.Context, productID int64) (item *domain.ProductOrderBlock, err error)
	CancelBlockForProduct(ctx context.Context, productID int64) (err error)
	ApplyBlockForProduct(ctx context.Context, productID int64) (err error)
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

func (uc *ProductOrderBlockInpl) CreateBlockForProduct(ctx context.Context, productID int64) (*domain.ProductOrderBlock, error) {
	return nil, nil
}

func (uc *ProductOrderBlockInpl) CancelBlockForProduct(ctx context.Context, productID int64) error {
	return nil
}

func (uc *ProductOrderBlockInpl) ApplyBlockForProduct(ctx context.Context, productID int64) error {
	return nil
}
