package postgres

import (
	"github.com/AbdulahadAbduqahhorov/gRPC/Ecommerce/order_service/genproto/order_service"
	"github.com/AbdulahadAbduqahhorov/gRPC/Ecommerce/order_service/storage/repo"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type orderRepo struct {
	db *sqlx.DB
}

func NewOrderRepo(db *sqlx.DB) repo.OrderRepoI {
	return &orderRepo{db: db}
}

func (r *orderRepo) CreateOrder(id string, req *order_service.CreateOrderRequest) (*order_service.CreateOrderResponse, error) {
	query := `INSERT INTO "order" (
				id,
				customer_name,
				customer_address,
				customer_phone
			) 
			VALUES ($1, $2, $3, $4) `

	_, err := r.db.Exec(query, id, req.CustomerName, req.CustomerAddress, req.CustomerPhone)

	if err != nil {
		return nil, err
	}

	query = `INSERT INTO order_items (
		id,
		order_id,
		product_id,
		quantity
	) 
	VALUES ($1, $2, $3, $4) `

	orderItems := req.Orderitems

	for key := range orderItems {
		t_id := uuid.New().String()
		_, err = r.db.Exec(query, t_id, id, orderItems[key].ProductId, orderItems[key].Quantity)

		if err != nil {
			return nil, err
		}
	}

	return &order_service.CreateOrderResponse{Id: id}, nil
}

func (r *orderRepo) GetOrderList(req *order_service.GetOrderListRequest) (*order_service.GetOrderListResponse, error) {
	res := &order_service.GetOrderListResponse{
		Orders: make([]*order_service.Order, 0),
	}
	query := `SELECT 
		id,
		customer_name,
		customer_address,
		customer_phone
	FROM "order"
	LIMIT $1
	OFFSET $2`

	rows, err := r.db.Queryx(query, int(req.Limit), int(req.Offset))
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var order order_service.Order
		err = rows.Scan(
			&order.Id,
			&order.CustomerName,
			&order.CustomerAddress,
			&order.CustomerPhone,
		)

		if err != nil {
			return nil, err
		}
		res.Orders = append(res.Orders, &order)
	}

	return res, nil
}

func (r *orderRepo) GetOrderById(id string) (*order_service.GetOrderByIdResponse, error) {
	res:=order_service.GetOrderByIdResponse{
		Orderitems: make([]*order_service.OrderItem, 0),
	}
	query := `
				SELECT
					id,
					customer_name,
					customer_address,
					customer_phone
				FROM
					"order"
				WHERE
					id = $1`
	row := r.db.QueryRow(query, id)
	err := row.Scan(
		&res.Id,
		&res.CustomerName,
		&res.CustomerAddress,
		&res.CustomerPhone,
	)
	if err != nil {
		return nil, err
	}

	queryOrderItems := `
				SELECT
					p.product_id,
					p.quantity
				FROM
					order_items p
				INNER JOIN
					"order" o
				ON
					p.order_id=o.id
				WHERE
					p.order_id = $1`

	rows, err:= r.db.Query(queryOrderItems, id)
	if err!=nil{
		return nil,err
	}
	for rows.Next() {
		var item order_service.OrderItem
		err = rows.Scan(
			&item.ProductId,
			&item.Quantity,
		)

		if err != nil {
			return nil, err
		}
		res.Orderitems=append(res.Orderitems, &item)
	}

	return &res, nil
}
