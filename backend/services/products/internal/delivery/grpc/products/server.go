package productsgrpc

import (
	"time"

	productsv1 "github.com/m11ano/mipt-webdev-course/backend/protos/gen/go/products"
	"github.com/m11ano/mipt-webdev-course/backend/services/products/internal/infra/config"
	"github.com/m11ano/mipt-webdev-course/backend/services/products/internal/usecase"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type serverAPI struct {
	productsv1.UnimplementedProductsServer
	cfg       config.Config
	productUC usecase.Product
}

func Register(gRPCServer *grpc.Server, cfg config.Config, productUC usecase.Product) {
	productsv1.RegisterProductsServer(gRPCServer, &serverAPI{
		productUC: productUC,
		cfg:       cfg,
	})
}

func toProtoTimestamp(t *time.Time) *timestamppb.Timestamp {
	if t == nil {
		return nil
	}
	return timestamppb.New(*t)
}
