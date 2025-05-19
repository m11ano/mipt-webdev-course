package productscl

import (
	"context"

	productsv1 "github.com/m11ano/mipt-webdev-course/backend/protos/gen/go/products"
)

func (c *ClientImpl) GetOrderBlockedProductsByOrderID(ctx context.Context, orderID int64) ([]*OrderBlockedProduct, error) {
	items, err := c.api.OrderBlockedProductsByOrderID(ctx, &productsv1.OrderBlockedProductsByOrderIDRequest{OrderId: orderID})
	if err != nil {
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
