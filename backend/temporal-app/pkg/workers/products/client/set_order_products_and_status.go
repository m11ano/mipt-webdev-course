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
var ErrSetOrderProductsAndStatusTimeout = e.NewErrorFrom(e.ErrBadRequest).SetMessage("req timeout")

type SetOrderProductsAndStatusFlowResult struct {
	we     tclient.WorkflowRun
	result workflows.SetOrderProductsAndStatusOut
	err    error
}

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

	execCtx := context.Background()

	startCh := make(chan SetOrderProductsAndStatusFlowResult, 1)

	go func() {
		we, err := c.client.ExecuteWorkflow(execCtx, options, workflows.SetOrderProductsAndStatus, workIn)
		if err != nil {
			startCh <- SetOrderProductsAndStatusFlowResult{err: e.NewErrorFrom(ErrWorkflowCantStart).Wrap(err)}
			return
		}

		if input.NotWait {
			startCh <- SetOrderProductsAndStatusFlowResult{we: we}
			return
		}

		var result workflows.SetOrderProductsAndStatusOut

		err = we.Get(execCtx, &result)
		if err != nil {
			startCh <- SetOrderProductsAndStatusFlowResult{err: e.NewErrorFrom(ErrWorkflowResutError).Wrap(err)}
			return
		}

		startCh <- SetOrderProductsAndStatusFlowResult{we: we, result: result}
	}()

	select {
	case <-ctx.Done():
		return e.NewErrorFrom(ErrSetOrderProductsAndStatusTimeout).Wrap(ctx.Err())
	case res := <-startCh:
		if res.err != nil {
			return res.err
		}

		if input.NotWait {
			return nil
		}

		if !res.result.IsOk {
			if res.result.ErrorCode == 1 {
				return ErrSetOrderProductsAndStatusCantReserve
			}

			if res.result.ErrorCode == 2 {
				return ErrSetOrderProductsAndStatusEmptyRequest
			}

			return e.NewErrorFrom(ErrWorkflowResutError)
		}

		return nil
	}
}
