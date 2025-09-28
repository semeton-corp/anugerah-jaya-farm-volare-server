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

.PHONY: build
build:
	@echo "Building the application..."
	go build -o app.exe cmd/app/main.go

.PHONY: clean-go-cache
clean-go-cache:
	@echo "Cleaning Go build cache..."
	go clean -cache -modcache

.PHONY: clean-docker-cache
clean-docker-cache:
	@echo "Cleaning Docker system..."
	docker system prune -af

.PHONY: compose-dev
compose-dev:
	@echo "Starting Docker containers in development mode..."
	docker compose -f compose.dev.yaml up --detach

.PHONY: compose-dev-down
compose-dev-down:
	@echo "Starting Docker containers in development mode..."
	docker compose -f compose.dev.yaml down

.PHONY: air
air:
	@echo "Starting go app with air..."
	air -c .air.toml