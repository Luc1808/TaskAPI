# --- Makefile ---

# Run app normally
run:
	go run ./cmd/api

# Run app with Air (hot reload)
dev:
	air -c ./air.toml

# Run all Go tests
test:
	go test ./...

# Run static analysis (if golangci-lint is installed)
lint:
	golangci-lint run

# Build binary to tmp/main
build:
	go build -o tmp/main ./cmd/api

migrate-up:
	set -a; . ./.env; set +a; \
	migrate -path migrations -database "$$DATABASE_URL" up

migrate-down:
	set -a; . ./.env; set +a; \
	migrate -path migrations -database "$$DATABASE_URL" down 1

migrate-force: ## if a migration gets stuck (DANGEROUS: set version manually)
	set -a; . ./.env; set +a; \
	migrate -path migrations -database "$$DATABASE_URL" force $(v)
	
migrate-version:
	set -a; . ./.env; set +a; \
	migrate -path migrations -database "$$DATABASE_URL" version

# usage: make migrate-new name=add_deadline_column
migrate-new:
	@[ -n "$(name)" ] || (echo "Usage: make migrate-new name=..." && exit 1)
	migrate create -ext sql -dir migrations -seq $(name)
