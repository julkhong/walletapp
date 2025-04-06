DB_CONTAINER_NAME=wallet-postgres
REDIS_CONTAINER_NAME=wallet-redis
SQL_SCHEMA_FILE=server/migrations/001_create_tables.sql
SQL_SEED_FILE=server/migrations/002_seed_dummy_users.sql

# Start Postgres container (if not already running)
db-up:
	@if docker ps -a --format '{{.Names}}' | grep -q "^$(DB_CONTAINER_NAME)$$"; then \
		echo "Container '$(DB_CONTAINER_NAME)' already exists. Skipping creation."; \
	else \
		echo "Starting Postgres container..."; \
		docker run --name $(DB_CONTAINER_NAME) \
			--env-file .env \
			-e POSTGRES_DB=$$(grep ^DB_NAME .env | cut -d '=' -f2) \
			-e POSTGRES_USER=$$(grep ^DB_USER .env | cut -d '=' -f2) \
			-e POSTGRES_PASSWORD=$$(grep ^DB_PASSWORD .env | cut -d '=' -f2) \
			-p $$(grep ^DB_PORT .env | cut -d '=' -f2):5432 \
			-d postgres || { echo "Failed to start Postgres container"; exit 1; }; \
	fi

# Start Redis container (if not already running)
redis-up:
	@if docker ps -a --format '{{.Names}}' | grep -q "^$(REDIS_CONTAINER_NAME)$$"; then \
		echo "Container '$(REDIS_CONTAINER_NAME)' already exists. Skipping creation."; \
	else \
		echo "Starting Redis container..."; \
		docker run --name $(REDIS_CONTAINER_NAME) \
			-p $$(grep ^REDIS_PORT .env | cut -d '=' -f2):6379 \
			-d redis || { echo "Failed to start Redis container"; exit 1; }; \
	fi

# Stop and remove Postgres container
db-down:
	@docker stop $(DB_CONTAINER_NAME) >/dev/null && docker rm -v $(DB_CONTAINER_NAME) >/dev/null

# Stop and remove Redis container
redis-down:
	@docker stop $(REDIS_CONTAINER_NAME) >/dev/null && docker rm -v $(REDIS_CONTAINER_NAME) >/dev/null

# Copy schema file into container
copy-schema:
	@docker cp $(SQL_SCHEMA_FILE) $(DB_CONTAINER_NAME):/tmp/schema.sql >/dev/null

# Copy seed file into container
copy-seed:
	@docker cp $(SQL_SEED_FILE) $(DB_CONTAINER_NAME):/tmp/schema.sql >/dev/null

# Wait for DB and run schema
db-init: copy-schema
	@sleep 5
	@docker exec -i $(DB_CONTAINER_NAME) \
	psql -U $$(grep ^DB_USER .env | cut -d '=' -f2) \
	-d $$(grep ^DB_NAME .env | cut -d '=' -f2) \
	-f /tmp/schema.sql >/dev/null || \
	{ echo "Schema execution failed"; exit 1; }

# Seed dummy data
db-seed: copy-seed
	@docker exec -i $(DB_CONTAINER_NAME) \
	psql -U $$(grep ^DB_USER .env | cut -d '=' -f2) \
	-d $$(grep ^DB_NAME .env | cut -d '=' -f2) \
	-f /tmp/schema.sql >/dev/null || \
	{ echo "Seeding failed"; exit 1; }

# Run Go app
run:
	@echo "Starting Go app..."
	@go run server/cmd/main.go || { echo "Go app failed to start"; exit 1; }

# One-liner to setup everything
start: db-up redis-up db-init db-seed run

# Full cleanup
clean:
	@docker stop $(DB_CONTAINER_NAME) >/dev/null && docker rm -v $(DB_CONTAINER_NAME) >/dev/null
	@docker stop $(REDIS_CONTAINER_NAME) >/dev/null && docker rm -v $(REDIS_CONTAINER_NAME) >/dev/null

# Run linter
lint:
	@golangci-lint run ./...

# Run linter with auto-fix (for formatting issues only)
lint-fix:
	@golangci-lint run --fix ./...

test:
	@echo "Running all tests with race detector and coverage..."
	@go test -race -coverprofile=coverage.out ./... -v

coverage:
	@go tool cover -html=coverage.out

.PHONY: db-up redis-up db-down redis-down copy-schema copy-seed
