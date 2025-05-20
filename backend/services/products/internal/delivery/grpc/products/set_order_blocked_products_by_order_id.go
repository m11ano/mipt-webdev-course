package productsgrpc

import (
	"context"
	"fmt"

	productsv1 "github.com/m11ano/mipt-webdev-course/backend/protos/gen/go/products"
)

func (s *serverAPI) SetOrderBlockedProductsByOrderID(ctx context.Context, in *productsv1.SetOrderBlockedProductsByOrderIDRequest) (*productsv1.SetOrderBlockedProductsByOrderIDResponse, error) {

	fmt.Println(in.GetOrderId())

	for _, item := range in.GetItems() {
		fmt.Println(item.GetProductId(), item.GetQuantity())
	}

	return nil, nil
}
