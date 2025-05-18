package usecase

import (
	"context"
	"log/slog"

	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/google/uuid"
	"github.com/m11ano/mipt-webdev-course/backend/services/products/internal/domain"
	"github.com/m11ano/mipt-webdev-course/backend/services/products/internal/infra/config"
	"github.com/m11ano/mipt-webdev-course/backend/services/products/internal/usecase/uctypes"
	"github.com/samber/lo"
)

type ProductSliderImageListOptions struct {
	ProductID *int64
	FileID    *uuid.UUID
}

//go:generate mockery --name=ProductSliderImage --output=../../tests/mocks --case=underscore
type ProductSliderImage interface {
	FindSliderImagesForProduct(ctx context.Context, productID int64) (items []*domain.ProductSliderImage, err error)
	SaveActualImagesForProductSlider(ctx context.Context, productID int64, filesIDs []uuid.UUID, markFilesAsAssigned bool) (items []*domain.ProductSliderImage, err error)
	DeleteImagesForProductSlider(ctx context.Context, productID int64) (err error)
}

//go:generate mockery --name=ProductSliderImageRepository --output=../../tests/mocks --case=underscore
type ProductSliderImageRepository interface {
	FindList(ctx context.Context, listOptions ProductSliderImageListOptions, queryParams *uctypes.QueryGetListParams) (items []*domain.ProductSliderImage, err error)
	Create(ctx context.Context, item *domain.ProductSliderImage) (err error)
	DeleteByList(ctx context.Context, listOptions ProductSliderImageListOptions) (err error)
}

type ProductSliderImageInpl struct {
	logger    *slog.Logger
	config    config.Config
	repo      ProductSliderImageRepository
	txManager *manager.Manager
	filesUC   File
}

func NewProductSliderImageInpl(logger *slog.Logger, config config.Config, txManager *manager.Manager, repo ProductSliderImageRepository, filesUC File) *ProductSliderImageInpl {
	uc := &ProductSliderImageInpl{
		logger:    logger,
		config:    config,
		txManager: txManager,
		repo:      repo,
		filesUC:   filesUC,
	}
	return uc
}

func (uc *ProductSliderImageInpl) FindSliderImagesForProduct(ctx context.Context, productID int64) ([]*domain.ProductSliderImage, error) {
	return uc.repo.FindList(ctx, ProductSliderImageListOptions{
		ProductID: &productID,
	}, nil)
}

func (uc *ProductSliderImageInpl) SaveActualImagesForProductSlider(ctx context.Context, productID int64, filesIDs []uuid.UUID, markFilesAsAssigned bool) ([]*domain.ProductSliderImage, error) {

	result := make([]*domain.ProductSliderImage, 0)

	err := uc.txManager.Do(ctx, func(ctx context.Context) error {
		curList, err := uc.repo.FindList(ctx, ProductSliderImageListOptions{
			ProductID: &productID,
		}, &uctypes.QueryGetListParams{
			ForUpdate: true,
		})
		if err != nil {
			return err
		}

		curListFilesIDs := make([]uuid.UUID, len(curList))
		for i, item := range curList {
			curListFilesIDs[i] = item.FileID
		}

		oldIDs, _ := lo.Difference(curListFilesIDs, filesIDs)

		err = uc.repo.DeleteByList(ctx, ProductSliderImageListOptions{
			ProductID: &productID,
		})
		if err != nil {
			return err
		}

		err = uc.filesUC.DeleteFilesByIDs(ctx, oldIDs)
		if err != nil {
			return err
		}

		sort := int32(-1)
		for _, fileID := range filesIDs {
			sort++
			item := domain.NewProductSliderImage(productID, fileID)
			item.Sort = sort

			result = append(result, item)

			err = uc.repo.Create(ctx, item)
			if err != nil {
				return err
			}
		}

		if markFilesAsAssigned {
			err = uc.filesUC.MarkAsAssigned(ctx, filesIDs)
			if err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (uc *ProductSliderImageInpl) DeleteImagesForProductSlider(ctx context.Context, productID int64) error {
	err := uc.txManager.Do(ctx, func(ctx context.Context) error {
		curList, err := uc.repo.FindList(ctx, ProductSliderImageListOptions{
			ProductID: &productID,
		}, &uctypes.QueryGetListParams{
			ForUpdate: true,
		})
		if err != nil {
			return err
		}

		curListFilesIDs := make([]uuid.UUID, len(curList))
		for i, item := range curList {
			curListFilesIDs[i] = item.FileID
		}

		err = uc.repo.DeleteByList(ctx, ProductSliderImageListOptions{
			ProductID: &productID,
		})
		if err != nil {
			return err
		}

		err = uc.filesUC.DeleteFilesByIDs(ctx, curListFilesIDs)
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
