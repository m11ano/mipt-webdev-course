package productstc

import (
	"context"
	"fmt"

	"github.com/m11ano/e"
	productsw "github.com/m11ano/mipt-webdev-course/backend/temporal-app/pkg/workers/products"
	"github.com/m11ano/mipt-webdev-course/backend/temporal-app/pkg/workers/products/workflows"
	tclient "go.temporal.io/sdk/client"
)

var ErrSetOrderProductsAndBlockCantReserve = e.NewErrorFrom(e.ErrBadRequest).SetMessage("can't reserve product")

func (c *ClientImpl) SetOrderProductsAndBlock(ctx context.Context, input SetOrderProductsAndBlockIn) error {
	options := tclient.StartWorkflowOptions{
		ID:        fmt.Sprintf("%s_%d", WorkflowOrderProductsPrefx, input.OrderID),
		TaskQueue: productsw.ProductsQueue,
	}

	workIn := workflows.BlockProductsAndConfirmIn{
		OrderID:       input.OrderID,
		OrderProducts: make([]workflows.OrderProductsItem, len(input.OrderProducts)),
	}

	for i, item := range input.OrderProducts {
		workIn.OrderProducts[i] = workflows.OrderProductsItem{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     item.Price,
		}
	}

	we, err := c.client.ExecuteWorkflow(ctx, options, workflows.SetOrderProductsAndBlock, workIn)
	if err != nil {
		return e.NewErrorFrom(ErrWorkflowCantStart).Wrap(err)
	}

	var result workflows.BlockProductsAndConfirmOut
	err = we.Get(ctx, &result)
	if err != nil {
		return e.NewErrorFrom(ErrWorkflowResutError).Wrap(err)
	}

	if !result.IsOk {
		if result.ErrorCode == 1 {
			return ErrSetOrderProductsAndBlockCantReserve
		}

		return e.NewErrorFrom(ErrWorkflowResutError).Wrap(err)
	}

	return nil
}
