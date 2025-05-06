include compose.env

.PHONY: run
run:
	@echo "Running the application..."
	go run cmd/app/main.go

.PHONY: compose-up
compose-up:
	@echo "Starting Docker containers..."
	docker compose --env-file compose.env up -d --build

.PHONY: compose-down
compose-down:
	@echo "Stopping Docker containers..."
	docker compose --env-file compose.env down

.PHONY: compose-debug
compose-debug:
	@echo "Starting Docker containers in debug mode..."
	docker compose --env-file compose.env up --build

.PHONY: build
build:
	@echo "Building the application..."
	go build -o app.exe cmd/app/main.go

get-db:
	@echo $(POSTGRES_DB)