DB_CONTAINER_NAME=wallet-postgres
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

# Stop and remove container
db-down:
	@docker stop $(DB_CONTAINER_NAME) >/dev/null && docker rm $(DB_CONTAINER_NAME) >/dev/null

# Copy schema file into container
copy-schema:
	@docker cp $(SQL_SCHEMA_FILE) $(DB_CONTAINER_NAME):/docker-entrypoint-initdb.d/schema.sql >/dev/null

# Copy seed file into container
copy-seed:
	@docker cp $(SQL_SEED_FILE) $(DB_CONTAINER_NAME):/docker-entrypoint-initdb.d/seed.sql >/dev/null

# Wait for DB and run schema
db-init: copy-schema
	@sleep 5
	@docker exec -i $(DB_CONTAINER_NAME) \
	psql -U $$(grep ^DB_USER .env | cut -d '=' -f2) \
	-d $$(grep ^DB_NAME .env | cut -d '=' -f2) \
	-f /docker-entrypoint-initdb.d/schema.sql >/dev/null || \
	{ echo "Schema execution failed"; exit 1; }

# Seed dummy data
db-seed: copy-seed
	@docker exec -i $(DB_CONTAINER_NAME) \
	psql -U $$(grep ^DB_USER .env | cut -d '=' -f2) \
	-d $$(grep ^DB_NAME .env | cut -d '=' -f2) \
	-f /docker-entrypoint-initdb.d/seed.sql >/dev/null || \
	{ echo "Seeding failed"; exit 1; }

# Run Go app
run:
	@echo "Starting Go app..."
	@go run server/cmd/main.go || { echo "Go app failed to start"; exit 1; }

# One-liner to setup everything
start: db-up db-init db-seed run

# Full cleanup
clean:
	@docker stop $(DB_CONTAINER_NAME) >/dev/null && docker rm $(DB_CONTAINER_NAME) >/dev/null
