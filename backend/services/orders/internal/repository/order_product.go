package repository

import (
	"context"
	"log/slog"
	"time"

	"github.com/Masterminds/squirrel"
	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/m11ano/e"
	"github.com/m11ano/mipt-webdev-course/backend/services/orders/internal/domain"
	"github.com/m11ano/mipt-webdev-course/backend/services/orders/internal/infra/db"
	"github.com/m11ano/mipt-webdev-course/backend/services/orders/internal/usecase"
	"github.com/m11ano/mipt-webdev-course/backend/services/orders/internal/usecase/uctypes"
	"github.com/m11ano/mipt-webdev-course/backend/services/orders/pkg/dbhelper"
	"github.com/shopspring/decimal"
)

const (
	orderProductTable = "order_product"
)

type DBOrderProduct struct {
	OrderID   int64           `db:"order_id"`
	ProductID int64           `db:"product_id"`
	Quantity  int32           `db:"quantity"`
	Price     decimal.Decimal `db:"price"`

	CreatedAt time.Time `db:"created_at"`
}

var (
	orderProductTableFields = []string{}
	orderProductDBSchema    = &DBOrderProduct{}
)

func init() {
	orderProductTableFields = dbhelper.ExtractDBFields(orderProductDBSchema)
}

type OrderProduct struct {
	logger *slog.Logger
	db     db.PgxPool
	txc    *trmpgx.CtxGetter
	qb     squirrel.StatementBuilderType
}

func NewOrderProduct(logger *slog.Logger, db db.PgxPool, txc *trmpgx.CtxGetter) *OrderProduct {
	return &OrderProduct{
		logger: logger,
		db:     db,
		txc:    txc,
		qb:     squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

func (r *OrderProduct) dbToDomain(db *DBOrderProduct) *domain.OrderProduct {
	return &domain.OrderProduct{
		ProductID: db.ProductID,
		OrderID:   db.OrderID,
		Quantity:  db.Quantity,
		Price:     db.Price,
		CreatedAt: db.CreatedAt,
	}
}

func (r *OrderProduct) buildWhereForList(listOptions usecase.OrderProductListOptions) squirrel.And {
	where := squirrel.And{}

	if listOptions.ProductID != nil {
		where = append(where, squirrel.Eq{"product_id": *listOptions.ProductID})
	}

	if listOptions.OrderID != nil {
		where = append(where, squirrel.Eq{"order_id": *listOptions.OrderID})
	}

	return where
}

func (r *OrderProduct) FindList(ctx context.Context, listOptions usecase.OrderProductListOptions, queryParams *uctypes.QueryGetListParams) ([]*domain.OrderProduct, error) {

	where := r.buildWhereForList(listOptions)

	q := r.qb.Select(orderProductTableFields...).From(orderProductTable).Where(where)

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

	dbData := []*DBOrderProduct{}

	if err := pgxscan.ScanAll(&dbData, rows); err != nil {
		errIsConv, convErr := e.ErrConvertPgxToLogic(err)
		if !errIsConv {
			r.logger.ErrorContext(ctx, "scan row", slog.Any("error", err))
		}
		return nil, convErr
	}

	result := make([]*domain.OrderProduct, 0, len(dbData))
	for _, dbItem := range dbData {
		result = append(result, r.dbToDomain(dbItem))
	}

	return result, nil
}

func (r *OrderProduct) Create(ctx context.Context, item *domain.OrderProduct) error {
	dataMap, err := dbhelper.StructToDBMap(item, orderProductDBSchema)
	if err != nil {
		r.logger.ErrorContext(ctx, "convert struct to db map", slog.Any("error", err))
		return e.NewErrorFrom(e.ErrInternal).Wrap(err)
	}

	query, args, err := r.qb.Insert(orderProductTable).SetMap(dataMap).ToSql()
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

func (r *OrderProduct) DeleteByList(ctx context.Context, listOptions usecase.OrderProductListOptions) error {

	where := r.buildWhereForList(listOptions)

	query, args, err := r.qb.Delete(orderProductTable).Where(where).ToSql()
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
