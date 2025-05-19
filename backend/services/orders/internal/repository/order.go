package repository

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/Masterminds/squirrel"
	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/google/uuid"
	"github.com/m11ano/e"
	"github.com/m11ano/mipt-webdev-course/backend/services/orders/internal/domain"
	"github.com/m11ano/mipt-webdev-course/backend/services/orders/internal/infra/db"
	"github.com/m11ano/mipt-webdev-course/backend/services/orders/internal/usecase"
	"github.com/m11ano/mipt-webdev-course/backend/services/orders/internal/usecase/uctypes"
	"github.com/m11ano/mipt-webdev-course/backend/services/orders/pkg/dbhelper"
	"github.com/shopspring/decimal"
	"golang.org/x/sync/errgroup"
)

const (
	orderTable = "order"
)

type DBOrder struct {
	ID              int64              `db:"id"`
	Status          domain.OrderStatus `db:"status"`
	OrderSum        decimal.Decimal    `db:"order_sum"`
	SecretKey       uuid.UUID          `db:"secret_key"`
	ClientName      string             `db:"client_name"`
	ClientSurname   string             `db:"client_surname"`
	ClientEmail     string             `db:"client_email"`
	ClientPhone     string             `db:"client_phone"`
	DeliveryAddress string             `db:"delivery_address"`

	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt *time.Time `db:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at"`
}

var (
	orderTableFields = []string{}
	orderDBSchema    = &DBOrder{}
)

func init() {
	orderTableFields = dbhelper.ExtractDBFields(orderDBSchema)
}

type Order struct {
	logger *slog.Logger
	db     db.PgxPool
	txc    *trmpgx.CtxGetter
	qb     squirrel.StatementBuilderType
}

func NewOrder(logger *slog.Logger, db db.PgxPool, txc *trmpgx.CtxGetter) *Order {
	return &Order{
		logger: logger,
		db:     db,
		txc:    txc,
		qb:     squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

func (r *Order) dbToDomain(db *DBOrder) *domain.Order {
	return &domain.Order{
		ID:              db.ID,
		Status:          db.Status,
		OrderSum:        db.OrderSum,
		SecretKey:       db.SecretKey,
		ClientName:      db.ClientName,
		ClientSurname:   db.ClientSurname,
		ClientEmail:     db.ClientEmail,
		ClientPhone:     db.ClientPhone,
		DeliveryAddress: db.DeliveryAddress,

		CreatedAt: db.CreatedAt,
		UpdatedAt: db.UpdatedAt,
		DeletedAt: db.DeletedAt,
	}
}

func (r *Order) buildWhereForList(listOptions usecase.OrderListOptions, withDeleted bool) squirrel.And {
	where := squirrel.And{}

	if listOptions.IDs != nil {
		where = append(where, squirrel.Eq{"id": *listOptions.IDs})
	}

	if !withDeleted {
		where = append(where, squirrel.Expr("deleted_at IS NULL"))
	}

	return where
}

var orderSortFieldMap = map[usecase.OrderListSortField]string{
	usecase.OrderListSortFieldID: "id",
}

func (r *Order) buildSortForList(listOptions usecase.OrderListOptions) []string {
	if listOptions.Sort == nil {
		return []string{}
	}

	sort := make([]string, 0, len(*listOptions.Sort))

	for _, sortItem := range *listOptions.Sort {
		sortField, ok := orderSortFieldMap[sortItem.Field]
		if ok {
			var dir string
			if sortItem.IsDesc {
				dir = "DESC"
			} else {
				dir = "ASC"
			}
			sort = append(sort, fmt.Sprintf("%s %s", sortField, dir))
		}
	}

	return sort
}

func (r *Order) buildPartUpdate(updateData usecase.OrderPartUpdateData) map[string]any {
	result := make(map[string]any)

	// if updateData.IsPublished != nil {
	// 	result["is_published"] = *updateData.IsPublished
	// }

	return result
}

func (r *Order) FindList(ctx context.Context, listOptions usecase.OrderListOptions, queryParams *uctypes.QueryGetListParams) ([]*domain.Order, error) {

	withDeleted := false
	if queryParams != nil && queryParams.WithDeleted {
		withDeleted = true
	}
	where := r.buildWhereForList(listOptions, withDeleted)
	sort := r.buildSortForList(listOptions)

	q := r.qb.Select(orderTableFields...).From(orderTable).Where(where).OrderBy(sort...)

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

	dbData := []*DBOrder{}

	if err := pgxscan.ScanAll(&dbData, rows); err != nil {
		errIsConv, convErr := e.ErrConvertPgxToLogic(err)
		if !errIsConv {
			r.logger.ErrorContext(ctx, "scan row", slog.Any("error", err))
		}
		return nil, convErr
	}

	result := make([]*domain.Order, 0, len(dbData))
	for _, dbItem := range dbData {
		result = append(result, r.dbToDomain(dbItem))
	}

	return result, nil
}

func (r *Order) FindPagedList(ctx context.Context, listOptions usecase.OrderListOptions, queryParams *uctypes.QueryGetListParams) ([]*domain.Order, int64, error) {

	withDeleted := false
	if queryParams != nil && queryParams.WithDeleted {
		withDeleted = true
	}
	where := r.buildWhereForList(listOptions, withDeleted)
	sort := r.buildSortForList(listOptions)

	q := r.qb.Select(orderTableFields...).From(orderTable).Where(where).OrderBy(sort...)
	qTotal := r.qb.Select("COUNT(*) as total").From(orderTable).Where(where)

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
		return nil, 0, e.NewErrorFrom(e.ErrInternal).Wrap(err)
	}

	queryTotal, argsTotal, err := qTotal.ToSql()
	if err != nil {
		r.logger.ErrorContext(ctx, "building total query", slog.Any("error", err))
		return nil, 0, e.NewErrorFrom(e.ErrInternal).Wrap(err)
	}

	var (
		dbData []*DBOrder
		total  int64
	)

	g, gCtx := errgroup.WithContext(ctx)

	g.Go(func() error {
		rows, err := r.txc.DefaultTrOrDB(gCtx, r.db).Query(gCtx, query, args...)
		if err != nil {
			errIsConv, convErr := e.ErrConvertPgxToLogic(err)
			if !errIsConv {
				r.logger.ErrorContext(ctx, "executing query", slog.Any("error", err))
			}
			return convErr
		}
		defer rows.Close()

		if err := pgxscan.ScanAll(&dbData, rows); err != nil {
			errIsConv, convErr := e.ErrConvertPgxToLogic(err)
			if !errIsConv {
				r.logger.ErrorContext(ctx, "scan row", slog.Any("error", err))
			}
			return convErr
		}
		return nil
	})

	g.Go(func() error {
		row := r.txc.DefaultTrOrDB(gCtx, r.db).QueryRow(gCtx, queryTotal, argsTotal...)
		if err := row.Scan(&total); err != nil {
			errIsConv, convErr := e.ErrConvertPgxToLogic(err)
			if !errIsConv {
				r.logger.ErrorContext(ctx, "scan total", slog.Any("error", err))
			}
			return convErr
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		return nil, 0, err
	}

	result := make([]*domain.Order, 0, len(dbData))
	for _, dbItem := range dbData {
		result = append(result, r.dbToDomain(dbItem))
	}

	return result, total, nil
}

func (r *Order) FindOneByID(ctx context.Context, id int64, queryParams *uctypes.QueryGetOneParams) (*domain.Order, error) {
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

	q := r.qb.Select(orderTableFields...).From(orderTable).Where(where)

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

	dbData := &DBOrder{}

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

func (r *Order) Create(ctx context.Context, item *domain.Order) error {
	dataMap, err := dbhelper.StructToDBMap(item, orderDBSchema)
	if err != nil {
		r.logger.ErrorContext(ctx, "convert struct to db map", slog.Any("error", err))
		return e.NewErrorFrom(e.ErrInternal).Wrap(err)
	}
	delete(dataMap, "id")
	delete(dataMap, "updated_at")
	delete(dataMap, "deleted_at")

	query, args, err := r.qb.Insert(orderTable).SetMap(dataMap).Suffix("RETURNING id").ToSql()
	if err != nil {
		r.logger.ErrorContext(ctx, "building query", slog.Any("error", err))
		return e.NewErrorFrom(e.ErrInternal).Wrap(err)
	}

	row := r.txc.DefaultTrOrDB(ctx, r.db).QueryRow(ctx, query, args...)

	if err := row.Scan(&item.ID); err != nil {
		errIsConv, convErr := e.ErrConvertPgxToLogic(err)
		if !errIsConv {
			r.logger.ErrorContext(ctx, "executing query", slog.Any("error", err))
		}
		return convErr
	}

	return nil
}

func (r *Order) Update(ctx context.Context, item *domain.Order) error {
	dataMap, err := dbhelper.StructToDBMap(item, orderDBSchema)
	if err != nil {
		r.logger.ErrorContext(ctx, "convert struct to db map", slog.Any("error", err))
		return e.NewErrorFrom(e.ErrInternal).Wrap(err)
	}
	delete(dataMap, "id")
	delete(dataMap, "created_at")
	delete(dataMap, "updated_at")
	delete(dataMap, "deleted_at")

	query, args, err := r.qb.Update(orderTable).Where(squirrel.Eq{"id": item.ID, "deleted_at": nil}).SetMap(dataMap).ToSql()
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

func (r *Order) PartUpdateByList(ctx context.Context, updateData usecase.OrderPartUpdateData, listOptions usecase.OrderListOptions, withDeleted bool) error {

	where := r.buildWhereForList(listOptions, withDeleted)

	dataMap := r.buildPartUpdate(updateData)

	query, args, err := r.qb.Update(orderTable).Where(where).SetMap(dataMap).ToSql()
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

func (r *Order) PartUpdateByID(ctx context.Context, updateData usecase.OrderPartUpdateData, id int64) error {

	dataMap := r.buildPartUpdate(updateData)

	query, args, err := r.qb.Update(orderTable).Where(squirrel.Eq{"id": id}).SetMap(dataMap).ToSql()
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

func (r *Order) DeleteByList(ctx context.Context, listOptions usecase.OrderListOptions) error {

	where := r.buildWhereForList(listOptions, false)

	dataMap := map[string]any{
		"deleted_at": time.Now(),
	}

	query, args, err := r.qb.Update(orderTable).Where(where).SetMap(dataMap).ToSql()
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
