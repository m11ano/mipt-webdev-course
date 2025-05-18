package usecase

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"runtime/debug"
	"time"

	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/google/uuid"
	"github.com/m11ano/mipt-webdev-course/backend/services/products/internal/domain"
	"github.com/m11ano/mipt-webdev-course/backend/services/products/internal/e"
	"github.com/m11ano/mipt-webdev-course/backend/services/products/internal/infra/config"
	"github.com/m11ano/mipt-webdev-course/backend/services/products/internal/infra/storage"
	"github.com/m11ano/mipt-webdev-course/backend/services/products/internal/usecase/uctypes"

	"github.com/disintegration/imaging"
	"github.com/samber/lo"
)

var ErrFileTargetNotSupport = e.NewErrorFrom(e.ErrInternal).SetMessage("file target not support")

type FilePartUpdateData struct {
	AssignedToTarget    *bool
	UploadedToStorage   *bool
	ToDeleteFromStorage *bool
}

type FileListOptions struct {
	IDs                 *[]uuid.UUID
	AssignedToTarget    *bool
	StorageFileKey      **string
	ToDeleteFromStorage *bool
	CreatedAt           *FileListOptionsCreatedAt
	DeletedAt           *FileListOptionsDeletedAt
}

type FileListOptionsCreatedAt struct {
	Value   *time.Time
	Compare uctypes.CompareType
}

type FileListOptionsDeletedAt struct {
	Value   *time.Time
	Compare uctypes.CompareType
}

//go:generate mockery --name=File --output=../../tests/mocks --case=underscore
type File interface {
	FindList(ctx context.Context, listOptions FileListOptions, queryParams *uctypes.QueryGetListParams) (items []*domain.File, err error)
	FindListInMap(ctx context.Context, listOptions FileListOptions, queryParams *uctypes.QueryGetListParams) (items map[uuid.UUID]*domain.File, err error)
	FindOneByID(ctx context.Context, id uuid.UUID, queryParams *uctypes.QueryGetOneParams) (file *domain.File, err error)
	MarkAsAssigned(ctx context.Context, ids []uuid.UUID) (err error)
	UploadImageFile(ctx context.Context, target domain.FileTarget, fileName string, fileData []byte) (file *domain.File, err error)
	DeleteFilesByIDs(ctx context.Context, ids []uuid.UUID) (err error)
	TaskRemoveFilesFromStorage(ctx context.Context, maxDeletedAt time.Time) (err error)
	TaskDeleteNotAssignedFiles(ctx context.Context, maxCreatedAt time.Time) (err error)
}

//go:generate mockery --name=FileRepository --output=../../tests/mocks --case=underscore
type FileRepository interface {
	FindList(ctx context.Context, listOptions FileListOptions, queryParams *uctypes.QueryGetListParams) (items []*domain.File, err error)
	FindOneByID(ctx context.Context, id uuid.UUID, queryParams *uctypes.QueryGetOneParams) (file *domain.File, err error)
	Create(ctx context.Context, item *domain.File) (err error)
	Update(ctx context.Context, item *domain.File) (err error)
	PartUpdateByList(ctx context.Context, updateData FilePartUpdateData, listOptions FileListOptions, withDeleted bool) (err error)
	PartUpdateByID(ctx context.Context, updateData FilePartUpdateData, id uuid.UUID) (err error)
	DeleteByList(ctx context.Context, listOptions FileListOptions, ToDeleteFromStorage bool) (err error)
}

type FileInpl struct {
	logger        *slog.Logger
	config        config.Config
	repo          FileRepository
	txManager     *manager.Manager
	storageClient storage.Client
}

func NewFileInpl(logger *slog.Logger, config config.Config, txManager *manager.Manager, repo FileRepository, storageClient storage.Client) *FileInpl {
	uc := &FileInpl{
		logger:        logger,
		config:        config,
		txManager:     txManager,
		repo:          repo,
		storageClient: storageClient,
	}
	return uc
}

var ImageSizes = map[domain.FileTarget][2]int{
	domain.FileTargetProductPreview: {500, 500},
	domain.FileTargetProductSlider:  {540, 540},
}

func (uc *FileInpl) resizeImg(ctx context.Context, target domain.FileTarget, fileData []byte) ([]byte, error) {

	sizes, ok := ImageSizes[target]
	if !ok {
		return nil, ErrFileTargetNotSupport
	}

	var reader io.Reader = bytes.NewReader(fileData)
	srcImg, err := imaging.Decode(reader)
	if err != nil {
		return nil, e.NewErrorFrom(e.ErrBadRequest).Wrap(err).SetMessage("invalid image file")
	}

	bounds := srcImg.Bounds()
	if bounds.Dx() < sizes[0] || bounds.Dy() < sizes[1] {
		return nil, e.NewErrorFrom(e.ErrBadRequest).SetMessage(fmt.Sprintf("min image resolution is %dx%d", sizes[0], sizes[1]))
	}

	resizedImg := imaging.Fill(
		srcImg,
		sizes[0],
		sizes[1],
		imaging.Center,
		imaging.Lanczos,
	)

	var resultBuf bytes.Buffer
	for _, q := range []int{95, 90, 85, 80, 75} {
		resultBuf.Reset()
		err := imaging.Encode(&resultBuf, resizedImg, imaging.JPEG, imaging.JPEGQuality(q))
		if err != nil {
			uc.logger.ErrorContext(ctx, "encoding image", slog.Any("error", err), slog.Any("quality", q), slog.Any("trackeback", string(debug.Stack())))
			return nil, e.NewErrorFrom(e.ErrInternal).Wrap(err).SetMessage("cannot encode image")
		}
		if resultBuf.Len() <= len(fileData) {
			break
		}
	}

	return resultBuf.Bytes(), nil
}

func (uc *FileInpl) FindList(ctx context.Context, listOptions FileListOptions, queryParams *uctypes.QueryGetListParams) ([]*domain.File, error) {
	return uc.repo.FindList(ctx, listOptions, queryParams)
}

func (uc *FileInpl) FindListInMap(ctx context.Context, listOptions FileListOptions, queryParams *uctypes.QueryGetListParams) (map[uuid.UUID]*domain.File, error) {
	list, err := uc.repo.FindList(ctx, listOptions, queryParams)
	if err != nil {

	}

	result := make(map[uuid.UUID]*domain.File, len(list))
	for _, item := range list {
		result[item.ID] = item
	}

	return result, nil
}

func (uc *FileInpl) FindOneByID(ctx context.Context, id uuid.UUID, queryParams *uctypes.QueryGetOneParams) (*domain.File, error) {
	return uc.repo.FindOneByID(ctx, id, queryParams)
}

func (uc *FileInpl) MarkAsAssigned(ctx context.Context, ids []uuid.UUID) error {
	return uc.repo.PartUpdateByList(ctx, FilePartUpdateData{
		AssignedToTarget: lo.ToPtr(true),
	}, FileListOptions{
		IDs: lo.ToPtr(ids),
	}, false)
}

func (uc *FileInpl) UploadImageFile(ctx context.Context, target domain.FileTarget, fileName string, fileData []byte) (*domain.File, error) {

	_, ok := ImageSizes[target]
	if !ok {
		return nil, ErrFileTargetNotSupport
	}

	resizeImg, err := uc.resizeImg(ctx, target, fileData)
	if err != nil {
		return nil, err
	}

	fileKey, _, _, err := uc.storageClient.UploadFileByBytes(ctx, storage.BucketProducts, fileName, resizeImg, nil)
	if err != nil {
		return nil, err
	}

	image := domain.NewFile(target, false)
	image.StorageFileKey = &fileKey
	image.UploadedToStorage = true

	err = uc.repo.Create(ctx, image)
	if err != nil {
		return nil, err
	}

	return image, nil
}

func (uc *FileInpl) DeleteFilesByIDs(ctx context.Context, ids []uuid.UUID) error {

	err := uc.repo.DeleteByList(ctx, FileListOptions{
		IDs: &ids,
	}, true)
	if err != nil {
		return err
	}

	return nil
}

func (uc *FileInpl) TaskRemoveFilesFromStorage(ctx context.Context, maxDeletedAt time.Time) error {

	deletedAt := &FileListOptionsDeletedAt{
		Value:   &maxDeletedAt,
		Compare: uctypes.CompareTypeLessOrEqual,
	}

	listOptions := FileListOptions{
		DeletedAt:           deletedAt,
		ToDeleteFromStorage: lo.ToPtr(true),
	}

	for {
		items, err := uc.repo.FindList(ctx, listOptions, &uctypes.QueryGetListParams{
			Limit:       1,
			WithDeleted: true,
		})
		if err != nil {
			return err
		}

		if len(items) == 0 {
			break
		}

		item := items[0]

		// Убрать из очереди на удаление если не нужно
		if item.StorageFileKey == nil || !item.UploadedToStorage {
			err := uc.repo.PartUpdateByID(ctx, FilePartUpdateData{
				ToDeleteFromStorage: lo.ToPtr(false),
				UploadedToStorage:   lo.ToPtr(false),
			}, item.ID)
			if err != nil {
				return err
			}

			continue
		}

		checkItems, err := uc.repo.FindList(ctx, FileListOptions{
			StorageFileKey: &item.StorageFileKey,
		}, &uctypes.QueryGetListParams{
			Limit:       1,
			WithDeleted: false,
		})
		if err != nil {
			return err
		}

		if len(checkItems) == 0 {
			// Удалить файл из хранилища если больше нет ссылок на него
			err = uc.storageClient.Delete(ctx, storage.BucketProducts, *item.StorageFileKey)
			if err != nil {
				return err
			}
		}

		err = uc.repo.PartUpdateByID(ctx, FilePartUpdateData{
			ToDeleteFromStorage: lo.ToPtr(false),
			UploadedToStorage:   lo.ToPtr(false),
		}, item.ID)
		if err != nil {
			return err
		}
	}

	return nil
}

func (uc *FileInpl) TaskDeleteNotAssignedFiles(ctx context.Context, maxCreatedAt time.Time) error {

	createdAt := &FileListOptionsCreatedAt{
		Value:   &maxCreatedAt,
		Compare: uctypes.CompareTypeLessOrEqual,
	}

	listOptions := FileListOptions{
		AssignedToTarget: lo.ToPtr(false),
		CreatedAt:        createdAt,
	}

	for {
		items, err := uc.repo.FindList(ctx, listOptions, &uctypes.QueryGetListParams{
			Limit: 100,
		})
		if err != nil {
			return err
		}

		if len(items) == 0 {
			break
		}

		IDs := make([]uuid.UUID, 0, len(items))
		for _, item := range items {
			IDs = append(IDs, item.ID)
		}

		err = uc.DeleteFilesByIDs(ctx, IDs)
		if err != nil {
			return err
		}
	}

	return nil
}
