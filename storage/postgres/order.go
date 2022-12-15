package postgres

import (
	"fmt"
	"strings"

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

func (r *orderRepo) CreateOrder(id string, total_price int, req *order_service.CreateOrderRequest) (*order_service.CreateOrderResponse, error) {
	query := `INSERT INTO "order" (
				id,
				user_id,
				customer_name,
				customer_address,
				customer_phone,
				total_price
			) 
			VALUES ($1, $2, $3, $4,$5,$6) `

	_, err := r.db.Exec(query, id, req.UserId, req.CustomerName, req.CustomerAddress, req.CustomerPhone, total_price)

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
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	if len(req.UserId) > 0 {
		setValues = append(setValues, fmt.Sprintf("AND user_id=$%d", argId))
		args = append(args, req.UserId)
		argId++
	}
	if req.Search != "" {
		setValues = append(setValues, fmt.Sprintf("AND (customer_name || customer_phone) ILIKE '%%' || $%d || '%%'", argId))
		args = append(args, req.Search)
		argId++
	}
	countQuery := `SELECT count(1) FROM "order"  WHERE true ` + strings.Join(setValues, " ")
	err := r.db.QueryRow(countQuery, args...).Scan(
		&res.Count,
	)
	if err != nil {
		return nil, err
	}
	if req.Limit > 0 {
		setValues = append(setValues, fmt.Sprintf("limit $%d ", argId))
		args = append(args, req.Limit)
		argId++
	}
	if req.Offset >= 0 {
		setValues = append(setValues, fmt.Sprintf("offset $%d ", argId))
		args = append(args, req.Offset)
		argId++
	}
	s := strings.Join(setValues, " ")
	query := `SELECT 
		id,
		user_id,
		customer_name,
		customer_address,
		customer_phone,
		total_price
	FROM "order" 
	WHERE true `+s

	rows, err := r.db.Queryx(query, args...)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var order order_service.Order
		err = rows.Scan(
			&order.Id,
			&order.UserId,
			&order.CustomerName,
			&order.CustomerAddress,
			&order.CustomerPhone,
			&order.TotalPrice,
		)

		if err != nil {
			return nil, err
		}
		res.Orders = append(res.Orders, &order)
	}

	return res, nil
}

func (r *orderRepo) GetOrderById(req *order_service.GetOrderByIdRequest) (*order_service.OrderInfo, error) {
	res := order_service.OrderInfo{
		Order:      &order_service.Order{},
		Orderitems: make([]*order_service.OrderItem, 0),
	}
	var filter string
	args:=make([]interface{},0)
	args=append(args, req.Id)
	if len(req.UserId)>0{
		filter="AND user_id=$2"
		args=append(args,req.UserId)
	}
	
	query := `
				SELECT
					id,
					user_id,
					customer_name,
					customer_address,
					customer_phone,
					total_price
				FROM
					"order"
				WHERE
					id = $1 `+filter
	row := r.db.QueryRow(query, args...)
	err := row.Scan(
		&res.Order.Id,
		&res.Order.UserId,
		&res.Order.CustomerName,
		&res.Order.CustomerAddress,
		&res.Order.CustomerPhone,
		&res.Order.TotalPrice,
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

	rows, err := r.db.Query(queryOrderItems, req.Id)
	if err != nil {
		return nil, err
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
		res.Orderitems = append(res.Orderitems, &item)
	}

	return &res, nil
}
