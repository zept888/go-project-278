.PHONY: help tidy dev build test

help:
	@echo "Available commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

tidy: ## Tidy up dependencies, format code, and run vet
	go mod tidy
	go fmt ./...
	go vet ./...

dev: ## Run the API server in development mode
	go run main.go

build: ## Build application binary
	go build -o bin/app main.go

test: ## Run all tests
	go test -v ./...
