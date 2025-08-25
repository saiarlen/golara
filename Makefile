# Golara Framework Makefile

.PHONY: build test test-unit test-integration clean install dev docs

# Build the application
build:
	go build -o bin/golara main.go

# Run all tests
test:
	go test ./tests/unit/... -v
	go test ./tests/integration/... -v
	go test ./tests/examples/... -v

# Run unit tests only
test-unit:
	go test ./tests/unit/... -v

# Run integration tests only
test-integration:
	go test ./tests/integration/... -v

# Run example tests only
test-examples:
	go test ./tests/examples/... -v

# Clean build artifacts
clean:
	rm -rf bin/
	rm -f golara

# Install dependencies
install:
	go mod tidy
	go mod download
	go install github.com/air-verse/air@latest

# Development server with hot reload
dev:
	@which air > /dev/null || (echo "Installing air..." && go install github.com/air-verse/air@latest)
	air

# Development server without hot reload
run:
	go run main.go

# Generate documentation
docs:
	@echo "Generating documentation..."
	@mkdir -p docs/html
	@go run scripts/generate_docs.go

# Database migrations
migrate:
	./bin/golara -subcommand migrate

# Rollback migrations
migrate-rollback:
	./bin/golara -subcommand migrate:rollback

# Generate controller
make-controller:
	@read -p "Controller name: " name; \
	./bin/golara -subcommand make:controller $$name

# Generate migration
make-migration:
	@read -p "Migration name: " name; \
	./bin/golara -subcommand make:migration $$name

# Setup project
setup: install
	cp .denv-example.yaml .denv.yaml
	mkdir -p storage/{logs,cache,app}
	mkdir -p bin
	@echo "✅ Project setup complete!"
	@echo "📝 Edit .denv.yaml with your configuration"
	@echo "🚀 Run 'make build && make migrate' to get started"