package orderscl

import "context"

type Client interface {
	SetOrderComposition(ctx context.Context, in SetOrderCompositionIn) (err error)
	CheckOrdersExistsByProductID(ctx context.Context, productID int64) (result bool, err error)
}
