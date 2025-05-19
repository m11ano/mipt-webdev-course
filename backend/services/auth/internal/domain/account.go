package domain

import (
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/m11ano/e"
	"golang.org/x/crypto/bcrypt"
)

var ErrAccountInvalidEmail = e.NewErrorFrom(e.ErrBadRequest).SetMessage("invalid email")

type Account struct {
	ID           uuid.UUID
	Name         string
	Surname      string
	Email        string
	PasswordHash string
	CreatedAt    time.Time
	UpdatedAt    *time.Time
	DeletedAt    *time.Time
}

func (a *Account) GeneretePasswordHash(password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		return err
	}

	a.PasswordHash = string(hash)

	return nil
}

func (a *Account) VerifyPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(a.PasswordHash), []byte(password))
	return err == nil
}

func (a *Account) SetEmail(email string) error {
	email = strings.ToLower(email)

	err := validate.Var(email, "required,email")
	if err != nil {
		return ErrAccountInvalidEmail
	}

	a.Email = email

	return nil
}

func NewAccount(name string, surname string, email string, password string) (*Account, error) {
	account := &Account{
		ID:        uuid.New(),
		Name:      name,
		Surname:   surname,
		CreatedAt: time.Now(),
	}

	err := account.SetEmail(email)
	if err != nil {
		return nil, err
	}

	err = account.GeneretePasswordHash(password)
	if err != nil {
		return nil, err
	}

	return account, nil
}
