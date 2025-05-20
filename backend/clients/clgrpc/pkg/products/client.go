package productscl

import "context"

type Client interface {
	GetProductsByIds(ctx context.Context, ids []int64) (items []*ProductListItem, err error)
	GetOrderBlockedProductsByOrderID(ctx context.Context, orderID int64) (items []*OrderBlockedProduct, err error)
	SetOrderBlockedProductsByOrderID(ctx context.Context, in SetOrderBlockedProductsByOrderIDIn) (err error)
}
