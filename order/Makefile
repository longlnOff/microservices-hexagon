include .env
export $(shell sed 's/=.*//' .env)

# Migration
db-status:
	goose status

up-all:
	goose up

up-by-one:
	goose up-by-one

up-to:
	goose up-to $(version)

down-by-one:
	goose down

down-to:
	goose down-to $(version)
	
migration-create:
	goose create -s $(name) sql
test-order-create-grpc:
	grpcurl -d '{"user_id": 123, "order_items": [{"product_code": "prod", "quantity": 4, "unit_price": 12}]}' -plaintext localhost:9000 Order/Create

# Database
run-postgres-db:
	docker run --name e_commercial_postgres -e POSTGRES_USER=longln -e POSTGRES_PASSWORD=mysecretpassword -p 5433:5432 -d postgres:16-alpine

create-db:
	docker exec -it e_commercial_postgres createdb --username=longln --owner=longln order

drop-db:
	docker exec -it e_commercial_postgres dropdb order

stop-db:
	docker stop e_commercial_postgres

start-db:
	docker start e_commercial_postgres

gen-docs:
	swag init -g ./main.go -d cmd,internal && swag fmt
