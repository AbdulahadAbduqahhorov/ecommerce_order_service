package service

import (
	"context"

	"github.com/AbdulahadAbduqahhorov/gRPC/Ecommerce/order_service/genproto/order_service"
	"github.com/AbdulahadAbduqahhorov/gRPC/Ecommerce/order_service/pkg/logger"
	"github.com/AbdulahadAbduqahhorov/gRPC/Ecommerce/order_service/storage"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type orderService struct {
	log logger.LoggerI
	storage storage.StorageI
	order_service.UnimplementedOrderServiceServer
}

func NewOrderService(log logger.LoggerI,db *sqlx.DB) *orderService {
	return &orderService{
		log:log,
		storage: storage.NewStoragePg(db),
	}
}

func (s *orderService) CreateOrder(ctx context.Context, req *order_service.CreateOrderRequest) (*order_service.CreateOrderResponse, error) {
	s.log.Info("---CreateOrder--->", logger.Any("req", req))
	id := uuid.New().String()
	res, err := s.storage.Order().CreateOrder(id, req)
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
