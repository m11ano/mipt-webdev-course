package productscl

import (
	"context"

	"github.com/m11ano/e"
	productsv1 "github.com/m11ano/mipt-webdev-course/backend/protos/gen/go/products"
)

func (c *ClientImpl) GetOrderBlockedProductsByOrderID(ctx context.Context, orderID int64) ([]*OrderBlockedProduct, error) {

	items, err := c.api.GetOrderBlockedProductsByOrderID(ctx, &productsv1.GetOrderBlockedProductsByOrderIDRequest{OrderId: orderID})
	if err != nil {
		if ok, lgErr := e.ErrConvertGRPCToLogic(err); ok {
			return nil, lgErr
		}
		return nil, err
	}

	result := make([]*OrderBlockedProduct, len(items.GetItems()))
	for i, item := range items.GetItems() {
		result[i] = &OrderBlockedProduct{
			ProductID: item.GetProductId(),
			Quantity:  item.GetQuantity(),
		}
	}

	return result, nil
}
