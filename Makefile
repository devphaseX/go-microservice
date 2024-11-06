FRONT_END_BINARY=frontApp
BROKER_BINARY=brokerApp
AUTH_BINARY=authApp
LOGGER_BINARY=loggerServiceApp

## up: starts all containers in the background without forcing build
up:
	@echo "Starting Docker images..."
	docker-compose up -d
	@echo "Docker images started!"

## up_build: stops docker-compose (if running), builds all projects and starts docker compose
up_build: build_broker build_auth build_logger
	@echo "Stopping docker images (if running...)"
	docker-compose down
	@echo "Building (when required) and starting docker images..."
	docker-compose up --build -d
	@echo "Docker images built and started!"

## down: stop docker compose
down:
	@echo "Stopping docker compose..."
	docker-compose down
	@echo "Done!"

## build_broker: builds the broker binary as a linux executable
build_broker:
	@echo "Building broker binary..."
	cd ./broker-service && env GOOS=linux CGO_ENABLED=0 go build -o ${BROKER_BINARY} ./cmd/api
	@echo "Done!"

## build_front: builds the auth binary
build_logger:
	@echo "Building logger binary..."
	cd ./logger-service && env CGO_ENABLED=0 go build -o ${LOGGER_BINARY} ./cmd/api
	@echo "Done!"


## build_front: builds the auth binary
build_auth:
	@echo "Building Auth service binary..."
	cd ./authentication-service && env CGO_ENABLED=0 go build -o ${AUTH_BINARY} ./cmd/api
	@echo "Done!"

## build_front: builds the frone end binary
build_front:
	@echo "Building front end binary..."
	cd ./front-end && env CGO_ENABLED=0 go build -o ${FRONT_END_BINARY} ./cmd/web
	@echo "Done!"


## start: starts the front end
start: build_front
	@echo "Starting front end"
	cd ./front-end && ./${FRONT_END_BINARY}

start_auth: build_auth
	@echo "starting auth server..."
	cd ./authenication-service && ./${AUTH_BINARY}

## stop: stop the front end
stop:
	@echo "Stopping front end..."
	@-pkill -SIGTERM -f "./${FRONT_END_BINARY}"
	@echo "Stopped front end!"


db_sqlc_auth:
	@echo "running sqlc"
	cd ./authenication-service && sqlc generate
	@echo "complete sqlc generation"

create_auth_migration:
	@echo "running auth migration"
	cd ./authenication-service &&  migrate create -ext sql -dir db/migrations -seq ${name}
	@echo "auth migration done"

createdb:
	docker exec -it postgres-container createdb --username=postgres --owner=postgres ${name}

dropdb:
	docker exec -it postgres-container dropdb --username=postgres ${name} -f

migrate_auth_db:
	 cd ./authenication-service &&	migrate -path db/migrations/ -database "postgresql://postgres:password@localhost:5432/service_auth?sslmode=disable" -verbose up


.PHONY: db_auth_generate
