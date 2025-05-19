package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/m11ano/e"
	"github.com/shopspring/decimal"
)

var ErrOrderCantSetStatus = e.NewErrorFrom(e.ErrBadRequest).SetMessage("cant set status")

type OrderStatus int

const (
	OrderStatusNew      OrderStatus = 0
	OrderStatusCreated  OrderStatus = 2
	OrderStatusInWork   OrderStatus = 3
	OrderStatusFinished OrderStatus = 10
	OrderStatusCanceled OrderStatus = 99
)

func (s OrderStatus) String() string {
	switch s {
	case OrderStatusNew:
		return "new"
	case OrderStatusCreated:
		return "created"
	case OrderStatusInWork:
		return "in_work"
	case OrderStatusFinished:
		return "finished"
	case OrderStatusCanceled:
		return "canceled"
	default:
		return ""
	}
}

var OrderStatusMap = map[string]OrderStatus{
	"new":      OrderStatusNew,
	"created":  OrderStatusCreated,
	"in_work":  OrderStatusInWork,
	"finished": OrderStatusFinished,
	"canceled": OrderStatusCanceled,
}

type Order struct {
	ID              int64
	Status          OrderStatus
	OrderSum        decimal.Decimal
	SecretKey       uuid.UUID
	ClientName      string
	ClientSurname   string
	ClientEmail     string
	ClientPhone     string
	DeliveryAddress string

	CreatedAt time.Time
	UpdatedAt *time.Time
	DeletedAt *time.Time
}

func NewOrder(id int64) *Order {
	return &Order{
		ID:        id,
		Status:    OrderStatusNew,
		SecretKey: uuid.New(),
		CreatedAt: time.Now(),
	}
}

func (p *Order) SetStatus(status OrderStatus) error {
	switch status {
	case OrderStatusNew:
		if p.Status != OrderStatusNew {
			return ErrOrderCantSetStatus
		}
	case OrderStatusCreated:
		if p.Status != OrderStatusNew {
			return ErrOrderCantSetStatus
		}
	case OrderStatusInWork:
		if p.Status != OrderStatusCreated {
			return ErrOrderCantSetStatus
		}
	case OrderStatusFinished:
		if p.Status != OrderStatusInWork {
			return ErrOrderCantSetStatus
		}
	case OrderStatusCanceled:
		if p.Status != OrderStatusFinished {
			return ErrOrderCantSetStatus
		}
	}

	p.Status = status

	return nil
}
