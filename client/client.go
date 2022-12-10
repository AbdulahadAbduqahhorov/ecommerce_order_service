package client

import (
	"github.com/AbdulahadAbduqahhorov/gRPC/Ecommerce/order_service/config"
	"github.com/AbdulahadAbduqahhorov/gRPC/Ecommerce/order_service/genproto/product_service"
	"google.golang.org/grpc"
)

type ServiceManagerI interface {
	ProductService() product_service.ProductServiceClient
}

type grpcClients struct {
	product product_service.ProductServiceClient
}

func NewGrpcClients(cfg config.Config) (ServiceManagerI, error) {
	connProductService, err := grpc.Dial(
		cfg.ProductServiceHost+cfg.ProductServicePort,
		grpc.WithInsecure(),
	)
	if err != nil {
		return nil, err
	}

	return &grpcClients{
		product: product_service.NewProductServiceClient(connProductService),
	}, nil
}

func (g *grpcClients) ProductService()  product_service.ProductServiceClient {
	return g.product
}
