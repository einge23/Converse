.PHONY: build run test lint clean

# Build the application
build:
	go build -o bin/converse main.go

# Run the application
run:
	go run main.go

# Run tests
test:
	go test -v ./...

# Run linters
lint:
	golangci-lint run

# Clean build artifacts
clean:
	rm -rf bin/
	go clean

# Install dependencies
deps:
	go mod tidy
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install golang.org/x/tools/cmd/goimports@latest

# Format code
fmt:
	go fmt ./...
	goimports -w . 