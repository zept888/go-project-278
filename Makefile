.PHONY: help tidy dev build test sqlc migrate

help:
	@echo "Available commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

tidy: ## Tidy up dependencies, format code, and run vet
	go mod tidy
	go fmt ./...
	go vet ./...

dev: ## Run API and frontend together (requires npm install)
	npm run dev

dev-api: ## Run the API server only
	go run .

build: ## Build application binary
	go build -o bin/app main.go

test: ## Run all tests
	go test -race -v ./...

sqlc: ## Generate data access code from SQL
	sqlc generate

migrate: ## Apply database migrations (requires DATABASE_URL)
	goose -dir ./db/migrations postgres "$(DATABASE_URL)" up
