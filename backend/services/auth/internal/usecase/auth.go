package usecase

import (
	"context"
	"errors"
	"log/slog"
	"strings"
	"time"

	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/golang-jwt/jwt/v5"
	"github.com/m11ano/e"
	"github.com/m11ano/mipt-webdev-course/backend/services/auth/internal/domain"
	"github.com/m11ano/mipt-webdev-course/backend/services/auth/internal/infra/config"
	"github.com/m11ano/mipt-webdev-course/backend/services/auth/pkg/auth"
)

//go:generate mockery --name=Auth --output=../../tests/mocks --case=underscore
type Auth interface {
	Login(ctx context.Context, email string, password string) (jwtToken string, account *domain.Account, err error)
}

type AuthInpl struct {
	logger         *slog.Logger
	config         config.Config
	txManager      *manager.Manager
	usecaseAccount Account
}

func NewAuthInpl(logger *slog.Logger, config config.Config, txManager *manager.Manager, usecaseAccount Account) *AuthInpl {
	uc := &AuthInpl{
		logger:         logger,
		config:         config,
		txManager:      txManager,
		usecaseAccount: usecaseAccount,
	}
	return uc
}

func (uc *AuthInpl) generateJWTToken(_ context.Context, claims *auth.AuthClaims) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(uc.config.Secrets.JWT))
}

func (uc *AuthInpl) Login(ctx context.Context, email string, password string) (string, *domain.Account, error) {
	email = strings.ToLower(email)

	account, err := uc.usecaseAccount.FindOneByEmail(ctx, email, nil)
	if err != nil {
		if errors.Is(err, e.ErrNotFound) {
			return "", nil, e.ErrUnauthorized
		}
		return "", nil, err
	}

	check := account.VerifyPassword(password)
	if !check {
		return "", nil, e.ErrUnauthorized
	}

	now := time.Now().UTC()

	claims := &auth.AuthClaims{
		AccountID: account.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Duration(uc.config.Auth.TokenTTLHours) * time.Hour)),
		},
	}

	token, err := uc.generateJWTToken(ctx, claims)
	if err != nil {
		return "", nil, err
	}

	return token, account, nil
}
