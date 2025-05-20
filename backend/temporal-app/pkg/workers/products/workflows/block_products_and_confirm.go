package workflows

import (
	"time"

	productscl "github.com/m11ano/mipt-webdev-course/backend/clients/clgrpc/pkg/products"
	"github.com/m11ano/mipt-webdev-course/backend/temporal-app/pkg/e2temperr"
	"github.com/m11ano/mipt-webdev-course/backend/temporal-app/pkg/workers/products/activities"
	"github.com/samber/lo"
	"github.com/shopspring/decimal"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

type OrderProductsItem struct {
	ProductID int64
	Quantity  int32
	Price     decimal.Decimal
}

type BlockProductsAndConfirmIn struct {
	OrderID       int64
	OrderProducts []OrderProductsItem
}

type BlockProductsAndConfirmOut struct {
	IsOk      bool
	ErrorCode int
}

func SetOrderProductsAndBlock(ctx workflow.Context, input BlockProductsAndConfirmIn) (*BlockProductsAndConfirmOut, error) {
	onceTryOpts := workflow.ActivityOptions{
		StartToCloseTimeout: time.Second * 1,
		RetryPolicy: &temporal.RetryPolicy{
			MaximumAttempts: 1,
		},
	}
	onceTryCtx := workflow.WithActivityOptions(ctx, onceTryOpts)

	//Бесконечный попытки пока не будет успеха или ответа 4xx
	unlimTryOpts := workflow.ActivityOptions{
		StartToCloseTimeout: time.Second * 1,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    time.Second * 1,
			BackoffCoefficient: 2.0,
			MaximumInterval:    time.Second * 30,
			MaximumAttempts:    0,
		},
	}
	unlimTryCtx := workflow.WithActivityOptions(ctx, unlimTryOpts)

	//Получим текущий список заблокированных товаров у заказа (если есть)
	currentOrderBlockedProducts := []*productscl.OrderBlockedProduct{}
	err := workflow.ExecuteActivity(onceTryCtx, "GetOrderBlockedProductsByOrderID", input.OrderID).Get(onceTryCtx, &currentOrderBlockedProducts)
	if err != nil {
		infBadInput := activities.InformOrdersServiceAboutOrderCompositionIn{
			OrderID: input.OrderID,
			IsOk:    false,
		}

		//Уведомим микросервис заказов о неуспешном блокировании товаров
		_ = workflow.ExecuteActivity(unlimTryCtx, "InformOrdersServiceAboutOrderComposition", infBadInput).Get(unlimTryCtx, nil)

		return &BlockProductsAndConfirmOut{
			IsOk:      false,
			ErrorCode: 99,
		}, nil
	}

	//Заблокируем список товаров для заказа
	blockInput := activities.SetOrderBlockedProductsByOrderIDIn{
		OrderID: input.OrderID,
		OrderProducts: lo.Map(input.OrderProducts, func(item OrderProductsItem, _ int) activities.SetOrderBlockedProductsByOrderIDItem {
			return activities.SetOrderBlockedProductsByOrderIDItem{
				ProductID: item.ProductID,
				Quantity:  item.Quantity,
			}
		}),
	}

	err = workflow.ExecuteActivity(onceTryCtx, "SetOrderBlockedProductsByOrderID", blockInput).Get(onceTryCtx, nil)
	if err != nil {
		//Если при блокировке истек таймаут или 500 - мы не знаем заблокировалось или нет, нужно отменять пока не будет успех
		needToCancel := false

		if ok, lgErr := e2temperr.TempErrConvertToLogicError(err); ok {
			if lgErr.Code() >= 500 && lgErr.Code() < 600 {
				needToCancel = true
			}
		} else if temporal.IsTimeoutError(err) {
			needToCancel = true
		}

		if needToCancel {
			//Возвращаем блокированные товары назад
			cancelBlockInput := activities.SetOrderBlockedProductsByOrderIDIn{
				OrderID: input.OrderID,
				OrderProducts: lo.Map(currentOrderBlockedProducts, func(item *productscl.OrderBlockedProduct, _ int) activities.SetOrderBlockedProductsByOrderIDItem {
					return activities.SetOrderBlockedProductsByOrderIDItem{
						ProductID: item.ProductID,
						Quantity:  item.Quantity,
					}
				}),
			}
			_ = workflow.ExecuteActivity(unlimTryCtx, "SetOrderBlockedProductsByOrderID", cancelBlockInput).Get(unlimTryCtx, nil)
		}

		infBadInput := activities.InformOrdersServiceAboutOrderCompositionIn{
			OrderID: input.OrderID,
			IsOk:    false,
		}

		//Уведомим микросервис заказов о неуспешном блокировании товаров
		_ = workflow.ExecuteActivity(unlimTryCtx, "InformOrdersServiceAboutOrderComposition", infBadInput).Get(unlimTryCtx, nil)

		return &BlockProductsAndConfirmOut{
			IsOk:      false,
			ErrorCode: 1,
		}, nil
	}

	//Уведомим микросервис заказов о успешной блокировке товаров
	infSuccessInput := activities.InformOrdersServiceAboutOrderCompositionIn{
		OrderID: input.OrderID,
		IsOk:    true,
		OrderProducts: lo.Map(input.OrderProducts, func(item OrderProductsItem, _ int) activities.InformOrdersServiceAboutOrderCompositionItem {
			return activities.InformOrdersServiceAboutOrderCompositionItem{
				ProductID: item.ProductID,
				Quantity:  item.Quantity,
				Price:     item.Price,
			}
		}),
	}

	err = workflow.ExecuteActivity(unlimTryCtx, "InformOrdersServiceAboutOrderComposition", infSuccessInput).Get(unlimTryCtx, nil)
	if err != nil {
		//Если ошибка - возвращаем блокированные товары назад
		cancelBlockInput := activities.SetOrderBlockedProductsByOrderIDIn{
			OrderID: input.OrderID,
			OrderProducts: lo.Map(currentOrderBlockedProducts, func(item *productscl.OrderBlockedProduct, _ int) activities.SetOrderBlockedProductsByOrderIDItem {
				return activities.SetOrderBlockedProductsByOrderIDItem{
					ProductID: item.ProductID,
					Quantity:  item.Quantity,
				}
			}),
		}
		_ = workflow.ExecuteActivity(unlimTryCtx, "SetOrderBlockedProductsByOrderID", cancelBlockInput).Get(unlimTryCtx, nil)

		infBadInput := activities.InformOrdersServiceAboutOrderCompositionIn{
			OrderID: input.OrderID,
			IsOk:    false,
		}

		_ = workflow.ExecuteActivity(unlimTryCtx, "InformOrdersServiceAboutOrderComposition", infBadInput).Get(unlimTryCtx, nil)

		return &BlockProductsAndConfirmOut{
			IsOk:      false,
			ErrorCode: 99,
		}, nil
	}

	out := &BlockProductsAndConfirmOut{
		IsOk:      true,
		ErrorCode: 0,
	}

	return out, nil
}
