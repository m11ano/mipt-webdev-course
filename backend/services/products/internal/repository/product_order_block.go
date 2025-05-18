package repository

import (
	"context"
	"log/slog"
	"time"

	"github.com/Masterminds/squirrel"
	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/m11ano/mipt-webdev-course/backend/services/products/internal/domain"
	"github.com/m11ano/mipt-webdev-course/backend/services/products/internal/e"
	"github.com/m11ano/mipt-webdev-course/backend/services/products/internal/infra/db"
	"github.com/m11ano/mipt-webdev-course/backend/services/products/internal/usecase"
	"github.com/m11ano/mipt-webdev-course/backend/services/products/internal/usecase/uctypes"
	"github.com/m11ano/mipt-webdev-course/backend/services/products/pkg/dbhelper"
)

const (
	productOrderBlockTable = "product_order_block"
)

type DBProductOrderBlock struct {
	ProductID int64 `db:"product_id"`
	OrderID   int64 `db:"order_id"`
	Quantity  int32 `db:"quantity"`

	CreatedAt time.Time `db:"created_at"`
}

var (
	productOrderBlockTableFields = []string{}
	productOrderBlockDBSchema    = &DBProductOrderBlock{}
)

func init() {
	productOrderBlockTableFields = dbhelper.ExtractDBFields(productOrderBlockDBSchema)
}

type ProductOrderBlock struct {
	logger *slog.Logger
	db     db.PgxPool
	txc    *trmpgx.CtxGetter
	qb     squirrel.StatementBuilderType
}

func NewProductOrderBlock(logger *slog.Logger, db db.PgxPool, txc *trmpgx.CtxGetter) *ProductOrderBlock {
	return &ProductOrderBlock{
		logger: logger,
		db:     db,
		txc:    txc,
		qb:     squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

func (r *ProductOrderBlock) dbToDomain(db *DBProductOrderBlock) *domain.ProductOrderBlock {
	return &domain.ProductOrderBlock{
		ProductID: db.ProductID,
		OrderID:   db.OrderID,
		Quantity:  db.Quantity,
		CreatedAt: db.CreatedAt,
	}
}

func (r *ProductOrderBlock) buildWhereForList(listOptions usecase.ProductOrderBlockListOptions) squirrel.And {
	where := squirrel.And{}

	if listOptions.ProductID != nil {
		where = append(where, squirrel.Eq{"product_id": *listOptions.ProductID})
	}

	if listOptions.OrderID != nil {
		where = append(where, squirrel.Eq{"order_id": *listOptions.OrderID})
	}

	return where
}

func (r *ProductOrderBlock) FindList(ctx context.Context, listOptions usecase.ProductOrderBlockListOptions, queryParams *uctypes.QueryGetListParams) ([]*domain.ProductOrderBlock, error) {

	where := r.buildWhereForList(listOptions)

	q := r.qb.Select(productOrderBlockTableFields...).From(productOrderBlockTable).Where(where)

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

	dbData := []*DBProductOrderBlock{}

	if err := pgxscan.ScanAll(&dbData, rows); err != nil {
		errIsConv, convErr := e.ErrConvertPgxToLogic(err)
		if !errIsConv {
			r.logger.ErrorContext(ctx, "scan row", slog.Any("error", err))
		}
		return nil, convErr
	}

	result := make([]*domain.ProductOrderBlock, 0, len(dbData))
	for _, dbItem := range dbData {
		result = append(result, r.dbToDomain(dbItem))
	}

	return result, nil
}

func (r *ProductOrderBlock) Create(ctx context.Context, item *domain.ProductOrderBlock) error {
	dataMap, err := dbhelper.StructToDBMap(item, productOrderBlockDBSchema)
	if err != nil {
		r.logger.ErrorContext(ctx, "convert struct to db map", slog.Any("error", err))
		return e.NewErrorFrom(e.ErrInternal).Wrap(err)
	}

	query, args, err := r.qb.Insert(productOrderBlockTable).SetMap(dataMap).ToSql()
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

func (r *ProductOrderBlock) DeleteByList(ctx context.Context, listOptions usecase.ProductOrderBlockListOptions) error {

	where := r.buildWhereForList(listOptions)

	query, args, err := r.qb.Delete(productOrderBlockTable).Where(where).ToSql()
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
