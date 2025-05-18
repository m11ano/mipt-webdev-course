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
	"github.com/m11ano/mipt-webdev-course/backend/services/products/internal/domain"
	"github.com/m11ano/mipt-webdev-course/backend/services/products/internal/e"
	"github.com/m11ano/mipt-webdev-course/backend/services/products/internal/infra/db"
	"github.com/m11ano/mipt-webdev-course/backend/services/products/internal/usecase"
	"github.com/m11ano/mipt-webdev-course/backend/services/products/internal/usecase/uctypes"
	"github.com/m11ano/mipt-webdev-course/backend/services/products/pkg/dbhelper"
	"github.com/shopspring/decimal"
	"golang.org/x/sync/errgroup"
)

const (
	productTable = "product"
)

type DBProduct struct {
	ID                 int64           `db:"id"`
	IsPublished        bool            `db:"is_published"`
	Name               string          `db:"name"`
	FullDescription    string          `db:"full_description"`
	Price              decimal.Decimal `db:"price"`
	StockAvailable     int32           `db:"stock_available"`
	ImagePreviewFileID *uuid.UUID      `db:"image_preview_file_id"`

	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt *time.Time `db:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at"`
}

var (
	productTableFields = []string{}
	productDBSchema    = &DBProduct{}
)

func init() {
	productTableFields = dbhelper.ExtractDBFields(productDBSchema)
}

type Product struct {
	logger *slog.Logger
	db     db.PgxPool
	txc    *trmpgx.CtxGetter
	qb     squirrel.StatementBuilderType
}

func NewProduct(logger *slog.Logger, db db.PgxPool, txc *trmpgx.CtxGetter) *Product {
	return &Product{
		logger: logger,
		db:     db,
		txc:    txc,
		qb:     squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

func (r *Product) dbToDomain(db *DBProduct) *domain.Product {
	return &domain.Product{
		ID:                 db.ID,
		IsPublished:        db.IsPublished,
		Name:               db.Name,
		FullDescription:    db.FullDescription,
		Price:              db.Price,
		StockAvailable:     db.StockAvailable,
		ImagePreviewFileID: db.ImagePreviewFileID,
		CreatedAt:          db.CreatedAt,
		UpdatedAt:          db.UpdatedAt,
		DeletedAt:          db.DeletedAt,
	}
}

func (r *Product) buildWhereForList(listOptions usecase.ProductListOptions, withDeleted bool) squirrel.And {
	where := squirrel.And{}

	if listOptions.IDs != nil {
		where = append(where, squirrel.Eq{"id": *listOptions.IDs})
	}

	if listOptions.IsPublished != nil {
		where = append(where, squirrel.Eq{"is_published": *listOptions.IsPublished})
	}

	if !withDeleted {
		where = append(where, squirrel.Expr("deleted_at IS NULL"))
	}

	return where
}

var productSortFieldMap = map[usecase.ProductListSortField]string{
	usecase.ProductListSortFieldCreatedAt: "created_at",
}

func (r *Product) buildSortForList(listOptions usecase.ProductListOptions) []string {
	if listOptions.Sort == nil {
		return []string{}
	}

	sort := make([]string, 0, len(*listOptions.Sort))

	for _, sortItem := range *listOptions.Sort {
		sortField, ok := productSortFieldMap[sortItem.Field]
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

func (r *Product) buildPartUpdate(updateData usecase.ProductPartUpdateData) map[string]any {
	result := make(map[string]any)

	if updateData.IsPublished != nil {
		result["is_published"] = *updateData.IsPublished
	}

	if updateData.Name != nil {
		result["name"] = *updateData.Name
	}

	if updateData.FullDescription != nil {
		result["full_description"] = *updateData.FullDescription
	}

	if updateData.Price != nil {
		result["price"] = *updateData.Price
	}

	if updateData.ImagePreviewID != nil {
		result["image_preview_file_id"] = *updateData.ImagePreviewID
	}

	if updateData.StockAvailable != nil {
		result["stock_available"] = *updateData.StockAvailable
	}

	return result
}

func (r *Product) FindList(ctx context.Context, listOptions usecase.ProductListOptions, queryParams *uctypes.QueryGetListParams) ([]*domain.Product, error) {

	withDeleted := false
	if queryParams != nil && queryParams.WithDeleted {
		withDeleted = true
	}
	where := r.buildWhereForList(listOptions, withDeleted)
	sort := r.buildSortForList(listOptions)

	q := r.qb.Select(productTableFields...).From(productTable).Where(where).OrderBy(sort...)

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

	dbData := []*DBProduct{}

	if err := pgxscan.ScanAll(&dbData, rows); err != nil {
		errIsConv, convErr := e.ErrConvertPgxToLogic(err)
		if !errIsConv {
			r.logger.ErrorContext(ctx, "scan row", slog.Any("error", err))
		}
		return nil, convErr
	}

	result := make([]*domain.Product, 0, len(dbData))
	for _, dbItem := range dbData {
		result = append(result, r.dbToDomain(dbItem))
	}

	return result, nil
}

func (r *Product) FindPagedList(ctx context.Context, listOptions usecase.ProductListOptions, queryParams *uctypes.QueryGetListParams) ([]*domain.Product, int64, error) {

	withDeleted := false
	if queryParams != nil && queryParams.WithDeleted {
		withDeleted = true
	}
	where := r.buildWhereForList(listOptions, withDeleted)
	sort := r.buildSortForList(listOptions)

	q := r.qb.Select(productTableFields...).From(productTable).Where(where).OrderBy(sort...)
	qTotal := r.qb.Select("COUNT(*) as total").From(productTable).Where(where)

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
		dbData []*DBProduct
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

	result := make([]*domain.Product, 0, len(dbData))
	for _, dbItem := range dbData {
		result = append(result, r.dbToDomain(dbItem))
	}

	return result, total, nil
}

func (r *Product) FindOneByID(ctx context.Context, id int64, queryParams *uctypes.QueryGetOneParams) (*domain.Product, error) {
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

	q := r.qb.Select(productTableFields...).From(productTable).Where(where)

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

	dbData := &DBProduct{}

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

func (r *Product) Create(ctx context.Context, item *domain.Product) error {
	dataMap, err := dbhelper.StructToDBMap(item, productDBSchema)
	if err != nil {
		r.logger.ErrorContext(ctx, "convert struct to db map", slog.Any("error", err))
		return e.NewErrorFrom(e.ErrInternal).Wrap(err)
	}
	delete(dataMap, "id")
	delete(dataMap, "updated_at")
	delete(dataMap, "deleted_at")

	query, args, err := r.qb.Insert(productTable).SetMap(dataMap).Suffix("RETURNING id").ToSql()
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

func (r *Product) Update(ctx context.Context, item *domain.Product) error {
	dataMap, err := dbhelper.StructToDBMap(item, productDBSchema)
	if err != nil {
		r.logger.ErrorContext(ctx, "convert struct to db map", slog.Any("error", err))
		return e.NewErrorFrom(e.ErrInternal).Wrap(err)
	}
	delete(dataMap, "id")
	delete(dataMap, "created_at")
	delete(dataMap, "updated_at")
	delete(dataMap, "deleted_at")

	query, args, err := r.qb.Update(productTable).Where(squirrel.Eq{"id": item.ID, "deleted_at": nil}).SetMap(dataMap).ToSql()
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

func (r *Product) PartUpdateByList(ctx context.Context, updateData usecase.ProductPartUpdateData, listOptions usecase.ProductListOptions, withDeleted bool) error {

	where := r.buildWhereForList(listOptions, withDeleted)

	dataMap := r.buildPartUpdate(updateData)

	query, args, err := r.qb.Update(productTable).Where(where).SetMap(dataMap).ToSql()
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

func (r *Product) PartUpdateByID(ctx context.Context, updateData usecase.ProductPartUpdateData, id int64) error {

	dataMap := r.buildPartUpdate(updateData)

	query, args, err := r.qb.Update(productTable).Where(squirrel.Eq{"id": id}).SetMap(dataMap).ToSql()
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

func (r *Product) DeleteByList(ctx context.Context, listOptions usecase.ProductListOptions) error {

	where := r.buildWhereForList(listOptions, false)

	dataMap := map[string]any{
		"deleted_at": time.Now(),
	}

	query, args, err := r.qb.Update(productTable).Where(where).SetMap(dataMap).ToSql()
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
