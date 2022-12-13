package repo

import "github.com/AbdulahadAbduqahhorov/gRPC/Ecommerce/order_service/genproto/order_service"

type OrderRepoI interface {
	CreateOrder(id string,total int, req *order_service.CreateOrderRequest) (*order_service.CreateOrderResponse, error)
	GetOrderList(req *order_service.GetOrderListRequest) (*order_service.GetOrderListResponse, error)
	GetOrderById(req *order_service.GetOrderByIdRequest) (*order_service.OrderInfo, error)
}
