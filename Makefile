.PHONY: gen clean

# Generate Swagger documentation
gen:
	cd backend && swag init -g cmd/server/main.go --parseDependency --parseInternal

# Clean generated files
clean:
	rm -rf backend/docs

# Install tools
tools:
	go install github.com/swaggo/swag/cmd/swag@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Run the server
run:
	cd backend && go run cmd/server/main.go

# Run tests
test:
	cd backend && go test ./...

# Format code
fmt:
	cd backend && go fmt ./...

# Lint code
lint:
	cd backend && golangci-lint run

# Build the server
build:
	cd backend && go build -o bin/server cmd/server/main.go 