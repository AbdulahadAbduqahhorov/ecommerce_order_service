package storage

import (
	"github.com/AbdulahadAbduqahhorov/gRPC/Ecommerce/order_service/storage/postgres"
	"github.com/AbdulahadAbduqahhorov/gRPC/Ecommerce/order_service/storage/repo"
	"github.com/jmoiron/sqlx"
)

type StorageI interface {
	Order() repo.OrderRepoI
}

type storagePg struct {
	db    *sqlx.DB
	order repo.OrderRepoI
}

func NewStoragePg(db *sqlx.DB) StorageI {
	return &storagePg{
		db:    db,
		order: postgres.NewOrderRepo(db),
	}
}

func (s *storagePg) Order() repo.OrderRepoI {
	return s.order
}
