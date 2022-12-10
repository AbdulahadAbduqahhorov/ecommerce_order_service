run: 
	go run cmd/main.go
migrateup:
	migrate -path ./migrations/postgres -database 'postgres://admin:admin@localhost:5432/order_db?sslmode=disable' up
migratedown:
	migrate -path ./migrations/postgres -database 'postgres://admin:admin@localhost:5432/order_db?sslmode=disable' down
pull-proto-module:
	git submodule update --init --recursive
update-proto-module:
	git submodule foreach --recursive git pull
