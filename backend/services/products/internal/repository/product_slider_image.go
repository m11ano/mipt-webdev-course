package repository

import (
	"context"
	"log/slog"
	"time"

	"github.com/Masterminds/squirrel"
	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/google/uuid"
	"github.com/m11ano/e"
	"github.com/m11ano/mipt-webdev-course/backend/services/products/internal/domain"
	"github.com/m11ano/mipt-webdev-course/backend/services/products/internal/infra/db"
	"github.com/m11ano/mipt-webdev-course/backend/services/products/internal/usecase"
	"github.com/m11ano/mipt-webdev-course/backend/services/products/internal/usecase/uctypes"
	"github.com/m11ano/mipt-webdev-course/backend/services/products/pkg/dbhelper"
)

const (
	productSliderImageTable = "product_slider_image"
)

type DBProductSliderImage struct {
	ProductID int64     `db:"product_id"`
	FileID    uuid.UUID `db:"file_id"`
	Sort      int32     `db:"sort"`

	CreatedAt time.Time `db:"created_at"`
}

var (
	productSliderImageTableFields = []string{}
	productSliderImageDBSchema    = &DBProductSliderImage{}
)

func init() {
	productSliderImageTableFields = dbhelper.ExtractDBFields(productSliderImageDBSchema)
}

type ProductSliderImage struct {
	logger *slog.Logger
	db     db.PgxPool
	txc    *trmpgx.CtxGetter
	qb     squirrel.StatementBuilderType
}

func NewProductSliderImage(logger *slog.Logger, db db.PgxPool, txc *trmpgx.CtxGetter) *ProductSliderImage {
	return &ProductSliderImage{
		logger: logger,
		db:     db,
		txc:    txc,
		qb:     squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

func (r *ProductSliderImage) dbToDomain(db *DBProductSliderImage) *domain.ProductSliderImage {
	return &domain.ProductSliderImage{
		ProductID: db.ProductID,
		FileID:    db.FileID,
		Sort:      db.Sort,
		CreatedAt: db.CreatedAt,
	}
}

func (r *ProductSliderImage) buildWhereForList(listOptions usecase.ProductSliderImageListOptions) squirrel.And {
	where := squirrel.And{}

	if listOptions.ProductID != nil {
		where = append(where, squirrel.Eq{"product_id": *listOptions.ProductID})
	}

	if listOptions.FileID != nil {
		where = append(where, squirrel.Eq{"file_id": *listOptions.FileID})
	}

	return where
}

func (r *ProductSliderImage) FindList(ctx context.Context, listOptions usecase.ProductSliderImageListOptions, queryParams *uctypes.QueryGetListParams) ([]*domain.ProductSliderImage, error) {

	where := r.buildWhereForList(listOptions)

	q := r.qb.Select(productSliderImageTableFields...).From(productSliderImageTable).Where(where).OrderBy("sort ASC")

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

	dbData := []*DBProductSliderImage{}

	if err := pgxscan.ScanAll(&dbData, rows); err != nil {
		errIsConv, convErr := e.ErrConvertPgxToLogic(err)
		if !errIsConv {
			r.logger.ErrorContext(ctx, "scan row", slog.Any("error", err))
		}
		return nil, convErr
	}

	result := make([]*domain.ProductSliderImage, 0, len(dbData))
	for _, dbItem := range dbData {
		result = append(result, r.dbToDomain(dbItem))
	}

	return result, nil
}

func (r *ProductSliderImage) Create(ctx context.Context, item *domain.ProductSliderImage) error {
	dataMap, err := dbhelper.StructToDBMap(item, productSliderImageDBSchema)
	if err != nil {
		r.logger.ErrorContext(ctx, "convert struct to db map", slog.Any("error", err))
		return e.NewErrorFrom(e.ErrInternal).Wrap(err)
	}

	query, args, err := r.qb.Insert(productSliderImageTable).SetMap(dataMap).ToSql()
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

func (r *ProductSliderImage) DeleteByList(ctx context.Context, listOptions usecase.ProductSliderImageListOptions) error {

	where := r.buildWhereForList(listOptions)

	query, args, err := r.qb.Delete(productSliderImageTable).Where(where).ToSql()
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
