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

# Database
run-postgres-db:
	docker run --name e_commercial_postgres -e POSTGRES_USER=longln -e POSTGRES_PASSWORD=mysecretpassword -p 5433:5432 -d postgres:16-alpine

create-db:
	docker exec -it e_commercial_postgres createdb --username=longln --owner=longln payment

drop-db:
	docker exec -it e_commercial_postgres dropdb payment

stop-db:
	docker stop e_commercial_postgres

start-db:
	docker start e_commercial_postgres

gen-docs:
	swag init -g ./main.go -d cmd,internal && swag fmt



run_docker_mysql:
	docker run --name mysql-microservice -p 3306:3306 -e MYSQL_ROOT_PASSWORD=verysecretpass -e MYSQL_DATABASE=order -d mysql 
create_db:
	docker exec mysql-microservice mysql -uroot -pverysecretpass -e "CREATE DATABASE payment; CREATE DATABASE shipping;"
start_docker_mysql:
	docker start mysql-microservice
stop_docker_mysql:
	docker stop mysql-microservice
delete_docker_mysql:
	docker rm mysql-microservice

service_payment:
	DB_DRIVER=mysql DATA_SOURCE_URL=root:verysecretpass@tcp\(127.0.0.1:3306\)/payment APPLICATION_PORT=3001 ENV=development go run cmd/main.go
