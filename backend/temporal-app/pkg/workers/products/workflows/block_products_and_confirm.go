package workflows

import (
	"time"

	productscl "github.com/m11ano/mipt-webdev-course/backend/clients/clgrpc/pkg/products"
	"go.temporal.io/sdk/workflow"
)

type OrderProductsItem struct {
	ProductID int64
	Quantity  int32
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

	opts := workflow.ActivityOptions{
		StartToCloseTimeout:    time.Second * 1,
		ScheduleToCloseTimeout: time.Second * 1,
	}
	ctx = workflow.WithActivityOptions(ctx, opts)

	currentOrderBlockedProducts := []*productscl.OrderBlockedProduct{}
	err := workflow.ExecuteActivity(ctx, "GetOrderBlockedProductsByOrderID", input.OrderID).Get(ctx, &currentOrderBlockedProducts)

	if err != nil {
		return &BlockProductsAndConfirmOut{
			IsOk:      false,
			ErrorCode: 99,
		}, err
	}

	//return nil, workflow.NewContinueAsNewError(ctx, SetOrderProductsAndBlock, input)

	//Заблочить товар

	//-Если ок, сообщить микросервису заказов что все ок (мб подтвердить заказ)

	//--Если ок, завершаем флоу успешно
	//--Если ошибка, отменить блок товара

	//-Если ошибка - сообщить микросервису заказов (мб удалить заказ если он новый)

	//--Если ок, завершаем флоу с ошибкой блокировки
	//--Если ошибка, отменить блок товара

	//---Если ок, завершаем флоу с ошибкой доступа к сервису заказов
	//---Если ошибка, ждем 3 секунды и перезапускаем флоу

	// Сообщить микросервису заказов что все ок (мб подтвердить заказ)

	//-Если ошибка, отменить блок товара

	out := &BlockProductsAndConfirmOut{
		IsOk:      false,
		ErrorCode: 0,
	}

	return out, nil
}
