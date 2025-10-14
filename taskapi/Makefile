# --- Makefile ---

# Run app normally
run:
	go run ./cmd/api

# Run app with Air (hot reload)
dev:
	air

# Run all Go tests
test:
	go test ./...

# Run static analysis (if golangci-lint is installed)
lint:
	golangci-lint run

# Build binary to tmp/main
build:
	go build -o tmp/main ./cmd/api

