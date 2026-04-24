.PHONY: create-migration migrate-up migrate-down migrate-drop docker-up docker-down docker-build docker-rebuild app simulate build clean logs

MIGRATE = migrate
MIGRATIONS_DIR = database/migrations
DB_URL = "mysql://root:password@tcp(127.0.0.1:3306)/network?charset=utf8mb4&parseTime=True&loc=Local"
DB_URL_TEST = "mysql://root:password@tcp(127.0.0.1:3306)/network_test?charset=utf8mb4&parseTime=True&loc=Local"

## Create new migration
create-migration:
	$(MIGRATE) create -ext sql -dir $(MIGRATIONS_DIR) $(name)

## Run migrations
migrate-up:
	$(MIGRATE) -database $(DB_URL) -path $(MIGRATIONS_DIR) up

## Setup and migrate test database
test-db:
	go run scripts/setup_test_db.go
	$(MIGRATE) -database $(DB_URL_TEST) -path $(MIGRATIONS_DIR) up

## Run all functional tests
test:
	go test -v ./test/...

## Rollback 1 step
migrate-down:
	$(MIGRATE) -database $(DB_URL) -path $(MIGRATIONS_DIR) down 1

## Drop all migrations
migrate-drop:
	$(MIGRATE) -database $(DB_URL) -path $(MIGRATIONS_DIR) drop -f

## Docker commands
docker-up:
	docker compose up -d

docker-down:
	docker compose down

app:
	go run cmd/app/main.go

simulate:
	go run cmd/simulate/main.go

## Build the application binary
build:
	go build -o bin/app ./cmd/app/main.go

## Build Docker images
docker-build:
	docker compose build

## Rebuild and restart services
docker-rebuild:
	docker compose up --build -d

## Show backend logs
logs:
	docker compose logs -f backend

## Clean up binaries and docker containers
clean: docker-down
	rm -rf bin/
