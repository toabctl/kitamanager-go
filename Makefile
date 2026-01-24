.PHONY: build run test test-unit test-integration test-contract test-fuzz test-coverage lint clean schema-docs docker-up docker-down docker-rebuild

# Build the application
build:
	go build -o bin/kitamanager-api ./cmd/api

# Run linter
lint:
	golangci-lint run ./...

# Run the application locally
run:
	go run ./cmd/api

# Run all tests (unit tests only, no tags)
test:
	go test -v ./...

# Run unit tests with race detection
test-unit:
	go test -v -race ./...

# Run integration tests (requires database)
test-integration:
	go test -v -race -tags=integration ./internal/integration/...

# Run API contract tests (requires database)
test-contract:
	go test -v -tags=contract ./internal/contract/...

# Run fuzz tests
test-fuzz:
	go test -fuzz=Fuzz -fuzztime=30s ./internal/models/...

# Run tests with coverage report
test-coverage:
	go test -v -race -coverprofile=coverage.out -covermode=atomic ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Clean build artifacts
clean:
	rm -rf bin/ coverage.out coverage.html

# Update database schema documentation (requires running database)
schema-docs:
	tbls doc --force

# Start docker containers
docker-up:
	docker compose up -d

# Stop docker containers
docker-down:
	docker compose down

# Rebuild and restart docker containers
docker-rebuild:
	docker compose up -d --build

# Run all CI checks locally
ci: lint test-unit build
	@echo "All CI checks passed!"
