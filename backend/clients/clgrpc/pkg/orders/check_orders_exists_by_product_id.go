package orderscl

import (
	"context"

	"github.com/m11ano/e"
	ordersv1 "github.com/m11ano/mipt-webdev-course/backend/protos/gen/go/orders"
)

func (c *ClientImpl) CheckOrdersExistsByProductID(ctx context.Context, productID int64) (bool, error) {

	result, err := c.api.CheckOrdersExistsByProductID(ctx, &ordersv1.CheckOrdersExistsByProductIDRequest{ProductId: productID})

	if err != nil {
		if ok, lgErr := e.ErrConvertGRPCToLogic(err); ok {
			return false, lgErr
		}

		return false, err
	}

	return result.Exists, nil
}
