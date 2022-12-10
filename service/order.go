package service

import (
	"context"
	"fmt"

	"github.com/AbdulahadAbduqahhorov/gRPC/Ecommerce/order_service/client"
	"github.com/AbdulahadAbduqahhorov/gRPC/Ecommerce/order_service/genproto/order_service"
	"github.com/AbdulahadAbduqahhorov/gRPC/Ecommerce/order_service/genproto/product_service"
	"github.com/AbdulahadAbduqahhorov/gRPC/Ecommerce/order_service/pkg/logger"
	"github.com/AbdulahadAbduqahhorov/gRPC/Ecommerce/order_service/storage"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type orderService struct {
	log      logger.LoggerI
	storage  storage.StorageI
	services client.ServiceManagerI
	order_service.UnimplementedOrderServiceServer
}

func NewOrderService(log logger.LoggerI, db *sqlx.DB, srvc client.ServiceManagerI) *orderService {
	return &orderService{
		log:      log,
		storage:  storage.NewStoragePg(db),
		services: srvc,
	}
}

func (s *orderService) CreateOrder(ctx context.Context, req *order_service.CreateOrderRequest) (*order_service.CreateOrderResponse, error) {
	s.log.Info("---CreateOrder--->", logger.Any("req", req))
	id := uuid.New().String()
	var total int
	for _, v := range req.Orderitems {
		product, err := s.services.ProductService().GetProductById(ctx, &product_service.GetProductByIdRequest{
			Id: v.ProductId,
		})
		if err != nil {
			s.log.Error("!!!CreateOrder--->", logger.Error(err))
			return nil, status.Error(codes.NotFound, fmt.Sprintf("Product Not found with id: %v", v.ProductId))
		}
		if v.Quantity > product.Quantity {
			s.log.Error("!!!CreateOrder--->", logger.Error(err))
			return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("Product with id:%v is available in %d units", v.ProductId,product.Quantity))
		}
	}
	for _, v := range req.Orderitems {
		product, _ := s.services.ProductService().GetProductById(ctx, &product_service.GetProductByIdRequest{
			Id: v.ProductId,
		})
		_, err:= s.services.ProductService().UpdateProduct(ctx, &product_service.UpdateProductRequest{
			Id:v.ProductId,
			Quantity: product.Quantity - v.Quantity,
		})
		if err != nil {
			s.log.Error("!!!CreateOrder--->", logger.Error(err))
			return nil, status.Error(codes.Internal, err.Error())
		}

		total += int(product.Price * v.Quantity)
	}

	res, err := s.storage.Order().CreateOrder(id, total, req)
	if err != nil {
		s.log.Error("!!!CreateOrder--->", logger.Error(err))
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return res, nil
}

func (s *orderService) GetOrderList(ctx context.Context, req *order_service.GetOrderListRequest) (*order_service.GetOrderListResponse, error) {
	s.log.Info("---GetOrderList--->", logger.Any("req", req))
	res, err := s.storage.Order().GetOrderList(req)
	if err != nil {
		s.log.Error("!!!GetOrderList--->", logger.Error(err))
		return res, status.Error(codes.Internal, err.Error())
	}

	return res, nil
}

func (s *orderService) GetOrderById(ctx context.Context, req *order_service.GetOrderByIdRequest) (*order_service.GetOrderByIdResponse, error) {
	s.log.Info("---GetOrderById--->", logger.Any("req", req))
	res, err := s.storage.Order().GetOrderById(req.Id)
	if err != nil {
		s.log.Error("!!!GetOrderById--->", logger.Error(err))
		return nil, status.Error(codes.NotFound, err.Error())
	}

	return res, nil
}
