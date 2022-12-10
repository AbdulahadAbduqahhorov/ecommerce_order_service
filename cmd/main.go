package main

import (
	"fmt"
	"net"

	"github.com/AbdulahadAbduqahhorov/gRPC/Ecommerce/order_service/client"
	"github.com/AbdulahadAbduqahhorov/gRPC/Ecommerce/order_service/config"
	"github.com/AbdulahadAbduqahhorov/gRPC/Ecommerce/order_service/genproto/order_service"
	"github.com/AbdulahadAbduqahhorov/gRPC/Ecommerce/order_service/pkg/logger"
	"github.com/AbdulahadAbduqahhorov/gRPC/Ecommerce/order_service/service"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"google.golang.org/grpc"
)

func main() {
	cfg := config.Load()
	log := logger.NewLogger("order_service", cfg.Environment)
	defer logger.Cleanup(log)
	conStr := fmt.Sprintf("host=%s port=%v user=%s password=%s dbname=%s sslmode=%s",
		cfg.PostgresHost,
		cfg.PostgresPort,
		cfg.PostgresUser,
		cfg.PostgresPassword,
		cfg.PostgresDatabase,
		"disable",
	)

	db, err := sqlx.Connect("postgres", conStr)
	if err != nil {
		log.Panic("postgres.NewPostgres", logger.Error(err))
	}

	lis, err := net.Listen("tcp", cfg.GrpcPort)
	if err != nil {
		log.Panic("net.Listen", logger.Error(err))
	}
	srvc,err:=client.NewGrpcClients(cfg)
	if err!=nil{
		log.Panic("client.NewGrpcClients()", logger.Error(err))

	}

	orderService := service.NewOrderService(log,db,srvc)
	s := grpc.NewServer()
	order_service.RegisterOrderServiceServer(s, orderService)
	
	log.Info("GRPC: Server being started...", logger.String("port", cfg.GrpcPort))

	if err := s.Serve(lis); err != nil {
		log.Panic("grpcServer.Serve", logger.Error(err))
	}
}
