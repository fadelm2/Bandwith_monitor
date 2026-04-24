APP_NAME=app
BUILD_DIR=bin
LOG_FILE=app.log
PID_FILE=app.pid
KILL_PORT=8082

MIGRATE=migrate
MIGRATIONS_DIR=database/migrations
DB_URL=mysql://root:password@tcp(127.0.0.1:3306)/network?charset=utf8mb4&parseTime=True&loc=Local
DB_URL_TEST=mysql://root:password@tcp(127.0.0.1:3306)/network_test?charset=utf8mb4&parseTime=True&loc=Local

########################################
# APP
########################################

## Build binary
build:
	mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(APP_NAME) cmd/app/main.go

## Run app (foreground)
run:
	go run cmd/app/main.go

## Kill process on port
kill-port:
	@echo "Checking port $(KILL_PORT)..."
	@PID=$$(lsof -t -i:$(KILL_PORT)); \
	if [ -n "$$PID" ]; then \
		echo "Killing PID $$PID on port $(KILL_PORT)"; \
		kill $$PID; \
	else \
		echo "No process running on port $(KILL_PORT)"; \
	fi

## Run app in background + log
run-bg: kill-port build
	nohup ./$(BUILD_DIR)/$(APP_NAME) > $(LOG_FILE) 2>&1 & echo $$! > $(PID_FILE)
	@echo "App running in background. PID: $$(cat $(PID_FILE))"

## Stop app
stop:
	@if [ -f $(PID_FILE) ]; then \
		kill $$(cat $(PID_FILE)) && rm -f $(PID_FILE); \
		echo "App stopped"; \
	else \
		echo "No PID file found"; \
	fi

## Restart app
restart: stop run-bg

## View logs
logs:
	tail -f $(LOG_FILE)

########################################
# MIGRATIONS
########################################

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

########################################
# TEST
########################################

## Setup and migrate test database
test-db:
	go run scripts/setup_test_db.go
	$(MIGRATE) -database $(DB_URL_TEST) -path $(MIGRATIONS_DIR) up

## Run all functional tests
test:
	go test -v ./test/...

########################################
# DOCKER
########################################

docker-up:
	docker compose up -d

docker-down:
	docker compose down
