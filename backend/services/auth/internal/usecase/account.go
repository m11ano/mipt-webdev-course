package usecase

import (
	"context"
	"log/slog"

	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/google/uuid"
	"github.com/m11ano/mipt-webdev-course/backend/services/auth/internal/domain"
	"github.com/m11ano/mipt-webdev-course/backend/services/auth/internal/infra/config"
	"github.com/m11ano/mipt-webdev-course/backend/services/auth/internal/usecase/uctypes"
)

//go:generate mockery --name=Account --output=../../tests/mocks --case=underscore
type Account interface {
	FindOneByEmail(ctx context.Context, email string, queryParams *uctypes.QueryGetOneParams) (account *domain.Account, err error)
	FindOneByID(ctx context.Context, id uuid.UUID, queryParams *uctypes.QueryGetOneParams) (account *domain.Account, err error)
}

//go:generate mockery --name=AccountRepository --output=../../tests/mocks --case=underscore
type AccountRepository interface {
	FindOneByEmail(ctx context.Context, email string, queryParams *uctypes.QueryGetOneParams) (account *domain.Account, err error)
	FindOneByID(ctx context.Context, id uuid.UUID, queryParams *uctypes.QueryGetOneParams) (account *domain.Account, err error)
}

type AccountInpl struct {
	logger    *slog.Logger
	config    config.Config
	repo      AccountRepository
	txManager *manager.Manager
}

func NewAccountInpl(logger *slog.Logger, config config.Config, txManager *manager.Manager, repo AccountRepository) *AccountInpl {
	uc := &AccountInpl{
		logger:    logger,
		config:    config,
		txManager: txManager,
		repo:      repo,
	}
	return uc
}

func (uc *AccountInpl) FindOneByEmail(ctx context.Context, email string, queryParams *uctypes.QueryGetOneParams) (*domain.Account, error) {
	return uc.repo.FindOneByEmail(ctx, email, queryParams)
}

func (uc *AccountInpl) FindOneByID(ctx context.Context, id uuid.UUID, queryParams *uctypes.QueryGetOneParams) (*domain.Account, error) {
	return uc.repo.FindOneByID(ctx, id, queryParams)
}
