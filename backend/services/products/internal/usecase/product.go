package usecase

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/google/uuid"
	"github.com/m11ano/mipt-webdev-course/backend/services/products/internal/domain"
	"github.com/m11ano/mipt-webdev-course/backend/services/products/internal/e"
	"github.com/m11ano/mipt-webdev-course/backend/services/products/internal/infra/config"
	"github.com/m11ano/mipt-webdev-course/backend/services/products/internal/usecase/uctypes"
	"github.com/samber/lo"
	"github.com/shopspring/decimal"
)

var ErrProductFileIDInvalid = e.NewErrorFrom(e.ErrBadRequest).SetMessage("invalid file_id")
var ErrProductAlreadyHasOrders = e.NewErrorFrom(e.ErrBadRequest).SetMessage("product already has orders")

type ProductPartUpdateData struct {
	Name            *string
	IsPublished     *bool
	FullDescription *string
	Price           *decimal.Decimal
	StockAvailable  *int32
	ImagePreviewID  **uuid.UUID
}

type ProductListSortField int

const (
	ProductListSortFieldCreatedAt ProductListSortField = iota
)

type ProductListSort struct {
	Field  ProductListSortField
	IsDesc bool
}

type ProductListOptions struct {
	IDs         *[]int64
	IsPublished *bool
	Sort        *[]ProductListSort
}

type ProductCreateIn struct {
	Product        *domain.Product
	SliderFilesIDs []uuid.UUID
}

type ProductUpdateIn struct {
	Name               string
	IsPublished        bool
	FullDescription    string
	Price              decimal.Decimal
	ImagePreviewFileID *uuid.UUID
	SliderFilesIDs     []uuid.UUID
}

type ProductOneFullOut struct {
	Product            *domain.Product
	ProductPreviewFile *domain.File
	SliderFiles        []*domain.File
}

type ProductFullOut struct {
	Product            *domain.Product
	ProductPreviewFile *domain.File
}

//go:generate mockery --name=Product --output=../../tests/mocks --case=underscore
type Product interface {
	FindFullPagedList(ctx context.Context, listOptions ProductListOptions, queryParams *uctypes.QueryGetListParams) (out []*ProductFullOut, total int64, err error)
	FindOneFullByID(ctx context.Context, id int64, queryParams *uctypes.QueryGetOneParams) (out *ProductOneFullOut, err error)
	Create(ctx context.Context, input ProductCreateIn) (product *domain.Product, slider []*domain.ProductSliderImage, err error)
	Update(ctx context.Context, id int64, input ProductUpdateIn) (product *domain.Product, slider []*domain.ProductSliderImage, err error)
	ChangeStock(ctx context.Context, id int64, value int32, isIncrease bool) (err error)
	Delete(ctx context.Context, id int64) (err error)
}

//go:generate mockery --name=ProductRepository --output=../../tests/mocks --case=underscore
type ProductRepository interface {
	FindList(ctx context.Context, listOptions ProductListOptions, queryParams *uctypes.QueryGetListParams) (items []*domain.Product, err error)
	FindPagedList(ctx context.Context, listOptions ProductListOptions, queryParams *uctypes.QueryGetListParams) (items []*domain.Product, total int64, err error)
	FindOneByID(ctx context.Context, id int64, queryParams *uctypes.QueryGetOneParams) (product *domain.Product, err error)
	Create(ctx context.Context, item *domain.Product) (err error)
	Update(ctx context.Context, item *domain.Product) (err error)
	PartUpdateByList(ctx context.Context, updateData ProductPartUpdateData, listOptions ProductListOptions, withDeleted bool) (err error)
	PartUpdateByID(ctx context.Context, updateData ProductPartUpdateData, id int64) (err error)
	DeleteByList(ctx context.Context, listOptions ProductListOptions) (err error)
}

type ProductInpl struct {
	logger               *slog.Logger
	config               config.Config
	repo                 ProductRepository
	txManager            *manager.Manager
	fileUC               File
	productSliderImageUC ProductSliderImage
	productOrderBlockUC  ProductOrderBlock
}

func NewProductInpl(logger *slog.Logger, config config.Config, txManager *manager.Manager, repo ProductRepository, filesUC File, productSliderImageUC ProductSliderImage, productOrderBlockUC ProductOrderBlock) *ProductInpl {
	uc := &ProductInpl{
		logger:               logger,
		config:               config,
		txManager:            txManager,
		repo:                 repo,
		fileUC:               filesUC,
		productSliderImageUC: productSliderImageUC,
		productOrderBlockUC:  productOrderBlockUC,
	}
	return uc
}

func (uc *ProductInpl) FindFullPagedList(ctx context.Context, listOptions ProductListOptions, queryParams *uctypes.QueryGetListParams) ([]*ProductFullOut, int64, error) {

	list, total, err := uc.repo.FindPagedList(ctx, listOptions, queryParams)
	if err != nil {
		return nil, 0, err
	}

	fileIDs := make([]uuid.UUID, 0, len(list))
	for _, item := range list {
		if item.ImagePreviewFileID != nil {
			fileIDs = append(fileIDs, *item.ImagePreviewFileID)
		}
	}

	files, err := uc.fileUC.FindListInMap(ctx, FileListOptions{
		IDs: &fileIDs,
	}, nil)
	if err != nil {
		return nil, 0, err
	}

	result := make([]*ProductFullOut, len(list))
	for i, item := range list {
		result[i] = &ProductFullOut{
			Product: item,
		}

		if item.ImagePreviewFileID != nil {
			file, ok := files[*item.ImagePreviewFileID]
			if ok {
				result[i].ProductPreviewFile = file
			}
		}
	}

	return result, total, nil
}

func (uc *ProductInpl) FindOneFullByID(ctx context.Context, id int64, queryParams *uctypes.QueryGetOneParams) (*ProductOneFullOut, error) {
	product, err := uc.repo.FindOneByID(ctx, id, queryParams)
	if err != nil {
		return nil, err
	}

	slider, err := uc.productSliderImageUC.FindSliderImagesForProduct(ctx, product.ID)
	if err != nil {
		return nil, err
	}

	filesIDsCap := len(slider)
	if product.ImagePreviewFileID != nil {
		filesIDsCap++
	}
	filesIDs := make([]uuid.UUID, 0, filesIDsCap)
	for _, item := range slider {
		filesIDs = append(filesIDs, item.FileID)
	}
	if product.ImagePreviewFileID != nil {
		filesIDs = append(filesIDs, *product.ImagePreviewFileID)
	}

	files, err := uc.fileUC.FindListInMap(ctx, FileListOptions{
		IDs: &filesIDs,
	}, nil)
	if err != nil {
		return nil, err
	}

	out := &ProductOneFullOut{
		Product:     product,
		SliderFiles: make([]*domain.File, 0, len(slider)),
	}

	if product.ImagePreviewFileID != nil {
		file, ok := files[*product.ImagePreviewFileID]
		if ok {
			out.ProductPreviewFile = file
		}
	}

	for _, item := range slider {
		file, ok := files[item.FileID]
		if ok {
			out.SliderFiles = append(out.SliderFiles, file)
		}
	}

	return out, nil
}

func (uc *ProductInpl) Create(ctx context.Context, input ProductCreateIn) (*domain.Product, []*domain.ProductSliderImage, error) {

	var resultSlider = make([]*domain.ProductSliderImage, 0)
	var err error

	input.SliderFilesIDs = lo.Uniq(input.SliderFilesIDs)

	filesIDs := make([]uuid.UUID, len(input.SliderFilesIDs))
	copy(filesIDs, input.SliderFilesIDs)
	if input.Product.ImagePreviewFileID != nil {
		filesIDs = append(filesIDs, *input.Product.ImagePreviewFileID)
	}

	err = uc.txManager.Do(ctx, func(ctx context.Context) error {

		files, err := uc.fileUC.FindListInMap(ctx, FileListOptions{
			IDs: &filesIDs,
		}, &uctypes.QueryGetListParams{
			ForUpdate: true,
		})
		if err != nil {
			return err
		}

		if input.Product.ImagePreviewFileID != nil {
			file, ok := files[*input.Product.ImagePreviewFileID]
			if !ok {
				return e.NewErrorFrom(ErrProductFileIDInvalid).SetMessage(fmt.Sprintf("incorrect file_id: %s", (*input.Product.ImagePreviewFileID).String()))
			}

			if file.Target != domain.FileTargetProductPreview {
				return e.NewErrorFrom(ErrProductFileIDInvalid).SetMessage(fmt.Sprintf("incorrect target, file_id: %s", file.ID.String()))
			}

			if file.AssignedToTarget {
				return e.NewErrorFrom(ErrProductFileIDInvalid).SetMessage(fmt.Sprintf("already assigned to target, file_id: %s", file.ID.String()))
			}
		}

		for _, fileID := range input.SliderFilesIDs {
			file, ok := files[fileID]
			if !ok {
				return e.NewErrorFrom(ErrProductFileIDInvalid).SetMessage(fmt.Sprintf("incorrect file_id: %s", fileID.String()))
			}

			if file.Target != domain.FileTargetProductSlider {
				return e.NewErrorFrom(ErrProductFileIDInvalid).SetMessage(fmt.Sprintf("incorrect target, file_id: %s", fileID.String()))
			}

			if file.AssignedToTarget {
				return e.NewErrorFrom(ErrProductFileIDInvalid).SetMessage(fmt.Sprintf("already assigned to target, file_id: %s", file.ID.String()))
			}
		}

		err = uc.repo.Create(ctx, input.Product)
		if err != nil {
			return err
		}

		err = uc.fileUC.MarkAsAssigned(ctx, filesIDs)
		if err != nil {
			return err
		}

		resultSlider, err = uc.productSliderImageUC.SaveActualImagesForProductSlider(ctx, input.Product.ID, input.SliderFilesIDs, false)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, nil, err
	}

	return input.Product, resultSlider, nil
}

func (uc *ProductInpl) Update(ctx context.Context, id int64, input ProductUpdateIn) (*domain.Product, []*domain.ProductSliderImage, error) {

	var product *domain.Product
	var resultSlider = make([]*domain.ProductSliderImage, 0)
	var err error

	input.SliderFilesIDs = lo.Uniq(input.SliderFilesIDs)

	newFilesIDs := make([]uuid.UUID, len(input.SliderFilesIDs))
	copy(newFilesIDs, input.SliderFilesIDs)
	if input.ImagePreviewFileID != nil {
		newFilesIDs = append(newFilesIDs, *input.ImagePreviewFileID)
	}

	err = uc.txManager.Do(ctx, func(ctx context.Context) error {

		product, err = uc.repo.FindOneByID(ctx, id, &uctypes.QueryGetOneParams{
			ForUpdate: true,
		})
		if err != nil {
			return err
		}

		curSliderImages, err := uc.productSliderImageUC.FindSliderImagesForProduct(ctx, product.ID)
		if err != nil {
			return err
		}

		curSliderImagesMap := make(map[uuid.UUID]*domain.ProductSliderImage, len(curSliderImages))
		for _, item := range curSliderImages {
			curSliderImagesMap[item.FileID] = item
		}

		newFiles, err := uc.fileUC.FindListInMap(ctx, FileListOptions{
			IDs: &newFilesIDs,
		}, &uctypes.QueryGetListParams{
			ForUpdate: true,
		})
		if err != nil {
			return err
		}

		if input.ImagePreviewFileID != nil {
			file, ok := newFiles[*input.ImagePreviewFileID]
			if !ok {
				return e.NewErrorFrom(ErrProductFileIDInvalid).SetMessage(fmt.Sprintf("incorrect file_id: %s", (*input.ImagePreviewFileID).String()))
			}

			if file.Target != domain.FileTargetProductPreview {
				return e.NewErrorFrom(ErrProductFileIDInvalid).SetMessage(fmt.Sprintf("incorrect target, file_id: %s", file.ID.String()))
			}

			if *input.ImagePreviewFileID != *product.ImagePreviewFileID && file.AssignedToTarget {
				return e.NewErrorFrom(ErrProductFileIDInvalid).SetMessage(fmt.Sprintf("already assigned to target, file_id: %s", file.ID.String()))
			}
		}

		var toDeleteOldPreview *uuid.UUID
		if input.ImagePreviewFileID != nil && *input.ImagePreviewFileID != *product.ImagePreviewFileID {
			toDeleteOldPreview = product.ImagePreviewFileID
		}

		for _, fileID := range input.SliderFilesIDs {
			file, ok := newFiles[fileID]
			if !ok {
				return e.NewErrorFrom(ErrProductFileIDInvalid).SetMessage(fmt.Sprintf("incorrect file_id: %s", fileID.String()))
			}

			if file.Target != domain.FileTargetProductSlider {
				return e.NewErrorFrom(ErrProductFileIDInvalid).SetMessage(fmt.Sprintf("incorrect target, file_id: %s", fileID.String()))
			}

			if _, ok = curSliderImagesMap[file.ID]; !ok && file.AssignedToTarget {
				return e.NewErrorFrom(ErrProductFileIDInvalid).SetMessage(fmt.Sprintf("already assigned to target, file_id: %s", file.ID.String()))
			}
		}

		product.Name = input.Name
		product.FullDescription = input.FullDescription
		product.IsPublished = input.IsPublished
		product.ImagePreviewFileID = input.ImagePreviewFileID
		err = product.SetPrice(input.Price)
		if err != nil {
			return err
		}

		err = uc.repo.Update(ctx, product)
		if err != nil {
			return err
		}

		err = uc.fileUC.MarkAsAssigned(ctx, newFilesIDs)
		if err != nil {
			return err
		}

		resultSlider, err = uc.productSliderImageUC.SaveActualImagesForProductSlider(ctx, product.ID, input.SliderFilesIDs, false)
		if err != nil {
			return err
		}

		if toDeleteOldPreview != nil {
			err = uc.fileUC.DeleteFilesByIDs(ctx, []uuid.UUID{*toDeleteOldPreview})
			if err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return nil, nil, err
	}

	return product, resultSlider, nil
}

func (uc *ProductInpl) ChangeStock(ctx context.Context, id int64, value int32, isIncrease bool) error {

	var product *domain.Product
	var err error

	err = uc.txManager.Do(ctx, func(ctx context.Context) error {

		product, err = uc.repo.FindOneByID(ctx, id, &uctypes.QueryGetOneParams{
			ForUpdate: true,
		})
		if err != nil {
			return err
		}

		if isIncrease {
			err = product.IncreaseStock(int64(value))
			if err != nil {
				return err
			}
		} else {
			err = product.DecreaseStock(int64(value))
			if err != nil {
				return err
			}
		}

		err = uc.repo.Update(ctx, product)
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

func (uc *ProductInpl) Delete(ctx context.Context, id int64) error {

	// 1) до транзакции проверим товар
	_, err := uc.repo.FindOneByID(ctx, id, nil)
	if err != nil {
		return err
	}

	// 2) TODO: до транзакции сделать проверку наличия заказов по товару, чтобы зря не создавать транзакцию и не блокировать товар, если он уже ранее был заказан

	// 3) создаем транзакцию
	err = uc.txManager.Do(ctx, func(ctx context.Context) error {

		// 4) блокируем товар внутри транзакции
		product, err := uc.repo.FindOneByID(ctx, id, &uctypes.QueryGetOneParams{
			ForUpdate: true,
		})
		if err != nil {
			return err
		}

		// 5) проверяем наличие блокировок на товар
		blockIsExists, err := uc.productOrderBlockUC.CheckBlockForProduct(ctx, product.ID)
		if err != nil {
			return err
		}

		if blockIsExists {
			return ErrProductAlreadyHasOrders
		}

		// 6) TODO: сделать проверку наличия заказов по товару

		// 7) удаляем товар
		err = uc.repo.DeleteByList(ctx, ProductListOptions{
			IDs: lo.ToPtr([]int64{product.ID}),
		})
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
