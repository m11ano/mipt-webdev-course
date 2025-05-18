package repository

import (
	"context"
	"log/slog"
	"time"

	"github.com/Masterminds/squirrel"
	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/google/uuid"
	"github.com/m11ano/mipt-webdev-course/backend/services/products/internal/domain"
	"github.com/m11ano/mipt-webdev-course/backend/services/products/internal/e"
	"github.com/m11ano/mipt-webdev-course/backend/services/products/internal/infra/db"
	"github.com/m11ano/mipt-webdev-course/backend/services/products/internal/usecase"
	"github.com/m11ano/mipt-webdev-course/backend/services/products/internal/usecase/uctypes"
	"github.com/m11ano/mipt-webdev-course/backend/services/products/pkg/dbhelper"
)

const (
	fileTable = "file"
)

type DBFile struct {
	ID                  uuid.UUID         `db:"id"`
	Target              domain.FileTarget `db:"file_target"`
	AssignedToTarget    bool              `db:"assigned_to_target"`
	StorageFileKey      *string           `db:"storage_file_key"`
	UploadedToStorage   bool              `db:"uploaded_to_storage"`
	ToDeleteFromStorage bool              `db:"to_delete_from_storage"`

	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt *time.Time `db:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at"`
}

var (
	fileTableFields = []string{}
	fileDBSchema    = &DBFile{}
)

func init() {
	fileTableFields = dbhelper.ExtractDBFields(fileDBSchema)
}

type File struct {
	logger *slog.Logger
	db     db.PgxPool
	txc    *trmpgx.CtxGetter
	qb     squirrel.StatementBuilderType
}

func NewFile(logger *slog.Logger, db db.PgxPool, txc *trmpgx.CtxGetter) *File {
	return &File{
		logger: logger,
		db:     db,
		txc:    txc,
		qb:     squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

func (r *File) dbToDomain(db *DBFile) *domain.File {
	return &domain.File{
		ID:                  db.ID,
		Target:              db.Target,
		AssignedToTarget:    db.AssignedToTarget,
		StorageFileKey:      db.StorageFileKey,
		UploadedToStorage:   db.UploadedToStorage,
		ToDeleteFromStorage: db.ToDeleteFromStorage,
		CreatedAt:           db.CreatedAt,
		UpdatedAt:           db.UpdatedAt,
		DeletedAt:           db.DeletedAt,
	}
}

func (r *File) buildWhereForList(listOptions usecase.FileListOptions, withDeleted bool) squirrel.And {
	where := squirrel.And{}

	if listOptions.IDs != nil {
		where = append(where, squirrel.Eq{"id": *listOptions.IDs})
	}

	if listOptions.StorageFileKey != nil {
		where = append(where, squirrel.Eq{"storage_file_key": *listOptions.StorageFileKey})
	}

	if listOptions.ToDeleteFromStorage != nil {
		where = append(where, squirrel.Eq{"to_delete_from_storage": *listOptions.ToDeleteFromStorage})
	}

	if listOptions.AssignedToTarget != nil {
		where = append(where, squirrel.Eq{"assigned_to_target": *listOptions.AssignedToTarget})
	}

	if listOptions.CreatedAt != nil {
		if listOptions.CreatedAt.Compare == uctypes.CompareTypeEqual {
			where = append(where, squirrel.Eq{"created_at": listOptions.CreatedAt.Value})
		} else {
			switch listOptions.CreatedAt.Compare {
			case uctypes.CompareTypeLess:
				where = append(where, squirrel.Expr("created_at < ?", *listOptions.CreatedAt.Value))
			case uctypes.CompareTypeLessOrEqual:
				where = append(where, squirrel.Expr("created_at <= ?", *listOptions.CreatedAt.Value))
			case uctypes.CompareTypeMore:
				where = append(where, squirrel.Expr("created_at > ?", *listOptions.CreatedAt.Value))
			case uctypes.CompareTypeMoreOrEqual:
				where = append(where, squirrel.Expr("created_at >= ?", *listOptions.CreatedAt.Value))
			}
		}
	}

	if listOptions.DeletedAt != nil {
		if listOptions.DeletedAt.Compare == uctypes.CompareTypeEqual {
			where = append(where, squirrel.Eq{"created_at": listOptions.DeletedAt.Value})
		} else {
			switch listOptions.DeletedAt.Compare {
			case uctypes.CompareTypeLess:
				where = append(where, squirrel.Expr("deleted_at < ?", *listOptions.DeletedAt.Value))
			case uctypes.CompareTypeLessOrEqual:
				where = append(where, squirrel.Expr("deleted_at <= ?", *listOptions.DeletedAt.Value))
			case uctypes.CompareTypeMore:
				where = append(where, squirrel.Expr("deleted_at > ?", *listOptions.DeletedAt.Value))
			case uctypes.CompareTypeMoreOrEqual:
				where = append(where, squirrel.Expr("deleted_at >= ?", *listOptions.DeletedAt.Value))
			}
		}
	}

	if !withDeleted {
		where = append(where, squirrel.Expr("deleted_at IS NULL"))
	}

	return where
}

func (r *File) buildPartUpdate(updateData usecase.FilePartUpdateData) map[string]any {
	result := make(map[string]any)

	if updateData.AssignedToTarget != nil {
		result["assigned_to_target"] = *updateData.AssignedToTarget
	}

	if updateData.UploadedToStorage != nil {
		result["uploaded_to_storage"] = *updateData.UploadedToStorage
	}

	if updateData.ToDeleteFromStorage != nil {
		result["to_delete_from_storage"] = *updateData.ToDeleteFromStorage
	}

	return result
}

func (r *File) FindList(ctx context.Context, listOptions usecase.FileListOptions, queryParams *uctypes.QueryGetListParams) ([]*domain.File, error) {

	withDeleted := false
	if queryParams != nil && queryParams.WithDeleted {
		withDeleted = true
	}
	where := r.buildWhereForList(listOptions, withDeleted)

	q := r.qb.Select(fileTableFields...).From(fileTable).Where(where)

	if queryParams != nil {
		if queryParams.ForUpdate {
			q = q.Suffix("FOR UPDATE")
		} else if queryParams.ForShare {
			q = q.Suffix("FOR SHARE")
		}

		if queryParams.Limit > 0 {
			q = q.Limit(queryParams.Limit)
		}

		if queryParams.Offset > 0 {
			q = q.Offset(queryParams.Offset)
		}
	}

	query, args, err := q.ToSql()
	if err != nil {
		r.logger.ErrorContext(ctx, "building query", slog.Any("error", err))
		return nil, e.NewErrorFrom(e.ErrInternal).Wrap(err)
	}

	rows, err := r.txc.DefaultTrOrDB(ctx, r.db).Query(ctx, query, args...)
	if err != nil {
		errIsConv, convErr := e.ErrConvertPgxToLogic(err)
		if !errIsConv {
			r.logger.ErrorContext(ctx, "executing query", slog.Any("error", err))
		}
		return nil, convErr
	}

	defer rows.Close()

	dbData := []*DBFile{}

	if err := pgxscan.ScanAll(&dbData, rows); err != nil {
		errIsConv, convErr := e.ErrConvertPgxToLogic(err)
		if !errIsConv {
			r.logger.ErrorContext(ctx, "scan row", slog.Any("error", err))
		}
		return nil, convErr
	}

	result := make([]*domain.File, 0, len(dbData))
	for _, dbItem := range dbData {
		result = append(result, r.dbToDomain(dbItem))
	}

	return result, nil
}

func (r *File) FindOneByID(ctx context.Context, id uuid.UUID, queryParams *uctypes.QueryGetOneParams) (*domain.File, error) {
	withDeleted := false
	if queryParams != nil && queryParams.WithDeleted {
		withDeleted = true
	}

	where := squirrel.And{
		squirrel.Eq{"id": id},
	}

	if !withDeleted {
		where = append(where, squirrel.Expr("deleted_at IS NULL"))
	}

	q := r.qb.Select(fileTableFields...).From(fileTable).Where(where)

	if queryParams != nil {
		if queryParams.ForUpdate {
			q = q.Suffix("FOR UPDATE")
		} else if queryParams.ForShare {
			q = q.Suffix("FOR SHARE")
		}
	}

	query, args, err := q.ToSql()
	if err != nil {
		r.logger.ErrorContext(ctx, "building query", slog.Any("error", err))
		return nil, e.NewErrorFrom(e.ErrInternal).Wrap(err)
	}

	rows, err := r.txc.DefaultTrOrDB(ctx, r.db).Query(ctx, query, args...)
	if err != nil {
		errIsConv, convErr := e.ErrConvertPgxToLogic(err)
		if !errIsConv {
			r.logger.ErrorContext(ctx, "executing query", slog.Any("error", err))
		}
		return nil, convErr
	}

	defer rows.Close()

	dbData := &DBFile{}

	if err := pgxscan.ScanOne(dbData, rows); err != nil {
		errIsConv, convErr := e.ErrConvertPgxToLogic(err)
		if !errIsConv {
			r.logger.ErrorContext(ctx, "scan row", slog.Any("error", err))
		}
		return nil, convErr
	}

	item := r.dbToDomain(dbData)

	return item, nil
}

func (r *File) Create(ctx context.Context, item *domain.File) error {
	dataMap, err := dbhelper.StructToDBMap(item, fileDBSchema)
	if err != nil {
		r.logger.ErrorContext(ctx, "convert struct to db map", slog.Any("error", err))
		return e.NewErrorFrom(e.ErrInternal).Wrap(err)
	}
	delete(dataMap, "updated_at")
	delete(dataMap, "deleted_at")

	query, args, err := r.qb.Insert(fileTable).SetMap(dataMap).ToSql()
	if err != nil {
		r.logger.ErrorContext(ctx, "building query", slog.Any("error", err))
		return e.NewErrorFrom(e.ErrInternal).Wrap(err)
	}

	_, err = r.txc.DefaultTrOrDB(ctx, r.db).Exec(ctx, query, args...)
	if err != nil {
		errIsConv, convErr := e.ErrConvertPgxToLogic(err)
		if !errIsConv {
			r.logger.ErrorContext(ctx, "executing query", slog.Any("error", err))
		}
		return convErr
	}

	return nil
}

func (r *File) Update(ctx context.Context, item *domain.File) error {
	dataMap, err := dbhelper.StructToDBMap(item, fileDBSchema)
	if err != nil {
		r.logger.ErrorContext(ctx, "convert struct to db map", slog.Any("error", err))
		return e.NewErrorFrom(e.ErrInternal).Wrap(err)
	}
	delete(dataMap, "id")
	delete(dataMap, "created_at")
	delete(dataMap, "updated_at")
	delete(dataMap, "deleted_at")

	query, args, err := r.qb.Update(fileTable).Where(squirrel.Eq{"id": item.ID, "deleted_at": nil}).SetMap(dataMap).ToSql()
	if err != nil {
		r.logger.ErrorContext(ctx, "building query", slog.Any("error", err))
		return e.NewErrorFrom(e.ErrInternal).Wrap(err)
	}

	_, err = r.txc.DefaultTrOrDB(ctx, r.db).Exec(ctx, query, args...)
	if err != nil {
		errIsConv, convErr := e.ErrConvertPgxToLogic(err)
		if !errIsConv {
			r.logger.ErrorContext(ctx, "executing query", slog.Any("error", err))
		}
		return convErr
	}

	return nil
}

func (r *File) PartUpdateByList(ctx context.Context, updateData usecase.FilePartUpdateData, listOptions usecase.FileListOptions, withDeleted bool) error {

	where := r.buildWhereForList(listOptions, withDeleted)

	dataMap := r.buildPartUpdate(updateData)

	query, args, err := r.qb.Update(fileTable).Where(where).SetMap(dataMap).ToSql()
	if err != nil {
		r.logger.ErrorContext(ctx, "building query", slog.Any("error", err))
		return e.NewErrorFrom(e.ErrInternal).Wrap(err)
	}

	_, err = r.txc.DefaultTrOrDB(ctx, r.db).Exec(ctx, query, args...)
	if err != nil {
		errIsConv, convErr := e.ErrConvertPgxToLogic(err)
		if !errIsConv {
			r.logger.ErrorContext(ctx, "executing query", slog.Any("error", err))
		}
		return convErr
	}
	return nil
}

func (r *File) PartUpdateByID(ctx context.Context, updateData usecase.FilePartUpdateData, id uuid.UUID) error {

	dataMap := r.buildPartUpdate(updateData)

	query, args, err := r.qb.Update(fileTable).Where(squirrel.Eq{"id": id}).SetMap(dataMap).ToSql()
	if err != nil {
		r.logger.ErrorContext(ctx, "building query", slog.Any("error", err))
		return e.NewErrorFrom(e.ErrInternal).Wrap(err)
	}

	_, err = r.txc.DefaultTrOrDB(ctx, r.db).Exec(ctx, query, args...)
	if err != nil {
		errIsConv, convErr := e.ErrConvertPgxToLogic(err)
		if !errIsConv {
			r.logger.ErrorContext(ctx, "executing query", slog.Any("error", err))
		}
		return convErr
	}
	return nil
}

func (r *File) DeleteByList(ctx context.Context, listOptions usecase.FileListOptions, ToDeleteFromStorage bool) error {

	where := r.buildWhereForList(listOptions, false)

	dataMap := map[string]any{
		"deleted_at": time.Now(),
	}

	if ToDeleteFromStorage {
		dataMap["to_delete_from_storage"] = true
	}

	query, args, err := r.qb.Update(fileTable).Where(where).SetMap(dataMap).ToSql()
	if err != nil {
		r.logger.ErrorContext(ctx, "building query", slog.Any("error", err))
		return e.NewErrorFrom(e.ErrInternal).Wrap(err)
	}

	_, err = r.txc.DefaultTrOrDB(ctx, r.db).Exec(ctx, query, args...)
	if err != nil {
		errIsConv, convErr := e.ErrConvertPgxToLogic(err)
		if !errIsConv {
			r.logger.ErrorContext(ctx, "executing query", slog.Any("error", err))
		}
		return convErr
	}
	return nil
}
