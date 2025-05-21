package productstc

import (
	"context"
	"fmt"

	"github.com/m11ano/e"
	productsw "github.com/m11ano/mipt-webdev-course/backend/temporal-app/pkg/workers/products"
	"github.com/m11ano/mipt-webdev-course/backend/temporal-app/pkg/workers/products/workflows"
	tclient "go.temporal.io/sdk/client"
)

var ErrSetOrderProductsAndStatusCantReserve = e.NewErrorFrom(e.ErrBadRequest).SetMessage("can't reserve product")
var ErrSetOrderProductsAndStatusEmptyRequest = e.NewErrorFrom(e.ErrBadRequest).SetMessage("empty request")

func (c *ClientImpl) SetOrderProductsAndStatus(ctx context.Context, input SetOrderProductsAndStatusIn) error {
	options := tclient.StartWorkflowOptions{
		ID:        fmt.Sprintf("%s_%d", WorkflowOrderProductsPrefx, input.OrderID),
		TaskQueue: productsw.ProductsQueue,
	}

	workIn := workflows.SetOrderProductsAndStatusIn{
		OrderID: input.OrderID,
	}

	if input.OrderProducts != nil {
		workInOrderProducts := make([]workflows.OrderProductsItem, len(*input.OrderProducts))
		for i, item := range *input.OrderProducts {
			workInOrderProducts[i] = workflows.OrderProductsItem{
				ProductID: item.ProductID,
				Quantity:  item.Quantity,
				Price:     item.Price,
			}
		}

		workIn.OrderProducts = &workInOrderProducts
	}

	if input.OrderStatus != nil {
		workIn.OrderStatus = input.OrderStatus
	}

	we, err := c.client.ExecuteWorkflow(ctx, options, workflows.SetOrderProductsAndStatus, workIn)
	if err != nil {
		return e.NewErrorFrom(ErrWorkflowCantStart).Wrap(err)
	}

	if input.NotWait {
		return nil
	}

	var result workflows.SetOrderProductsAndStatusOut
	err = we.Get(ctx, &result)
	if err != nil {
		return e.NewErrorFrom(ErrWorkflowResutError).Wrap(err)
	}

	if !result.IsOk {
		if result.ErrorCode == 1 {
			return ErrSetOrderProductsAndStatusCantReserve
		}

		if result.ErrorCode == 2 {
			return ErrSetOrderProductsAndStatusEmptyRequest
		}

		return e.NewErrorFrom(ErrWorkflowResutError).Wrap(err)
	}

	return nil
}
