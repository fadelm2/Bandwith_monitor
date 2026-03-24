.PHONY: create-migration migrate-up migrate-down migrate-drop docker-up docker-down app simulate

MIGRATE = migrate
MIGRATIONS_DIR = database/migrations
DB_URL = "mysql://root:password@tcp(127.0.0.1:3306)/network?charset=utf8mb4&parseTime=True&loc=Local"

## Create new migration
create-migration:
	$(MIGRATE) create -ext sql -dir $(MIGRATIONS_DIR) $(name)

## Run migrations
migrate-up:
	$(MIGRATE) -database $(DB_URL) -path $(MIGRATIONS_DIR) up

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
