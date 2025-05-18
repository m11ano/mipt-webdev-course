package repository

import (
	"context"
	"log/slog"
	"time"

	"github.com/Masterminds/squirrel"
	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/google/uuid"
	"github.com/m11ano/mipt-webdev-course/backend/services/auth/internal/domain"
	"github.com/m11ano/mipt-webdev-course/backend/services/auth/internal/e"
	"github.com/m11ano/mipt-webdev-course/backend/services/auth/internal/infra/db"
	"github.com/m11ano/mipt-webdev-course/backend/services/auth/internal/usecase/uctypes"
	"github.com/m11ano/mipt-webdev-course/backend/services/auth/pkg/dbhelper"
)

const (
	accountTable = "account"
)

type DBAccount struct {
	ID           uuid.UUID  `db:"id"`
	Name         string     `db:"name"`
	Surname      string     `db:"surname"`
	Email        string     `db:"email"`
	PasswordHash string     `db:"password_hash"`
	CreatedAt    time.Time  `db:"created_at"`
	UpdatedAt    *time.Time `db:"updated_at"`
}

var (
	accountTableFields = []string{}
	accountDBSchema    = &DBAccount{}
)

func init() {
	accountTableFields = dbhelper.ExtractDBFields(accountDBSchema)
}

type Account struct {
	logger *slog.Logger
	db     db.PgxPool
	txc    *trmpgx.CtxGetter
	qb     squirrel.StatementBuilderType
}

func NewAccount(logger *slog.Logger, db db.PgxPool, txc *trmpgx.CtxGetter) *Account {
	return &Account{
		logger: logger,
		db:     db,
		txc:    txc,
		qb:     squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

func (r *Account) dbToDomain(db *DBAccount) *domain.Account {
	return &domain.Account{
		ID:           db.ID,
		Name:         db.Name,
		Surname:      db.Surname,
		Email:        db.Email,
		PasswordHash: db.PasswordHash,
		CreatedAt:    db.CreatedAt,
		UpdatedAt:    db.UpdatedAt,
	}
}

func (r *Account) FindOneByEmail(ctx context.Context, email string, queryParams *uctypes.QueryGetOneParams) (*domain.Account, error) {
	withDeleted := false
	if queryParams != nil && queryParams.WithDeleted {
		withDeleted = true
	}

	where := squirrel.And{
		squirrel.Eq{"email": email},
	}

	if !withDeleted {
		where = append(where, squirrel.Expr("deleted_at IS NULL"))
	}

	q := r.qb.Select(accountTableFields...).From(accountTable).Where(where)

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

	dbData := &DBAccount{}

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

func (r *Account) FindOneByID(ctx context.Context, id uuid.UUID, queryParams *uctypes.QueryGetOneParams) (*domain.Account, error) {
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

	q := r.qb.Select(accountTableFields...).From(accountTable).Where(where)

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

	dbData := &DBAccount{}

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
