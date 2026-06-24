.PHONY: build run test test-cover generate migrate-up migrate-down lint tidy

APP=tps-pkb
BINARY=./bin/$(APP)
MIGRATIONS_PATH=migrations
DB_URL?=$(shell grep DATABASE_URL .env 2>/dev/null | cut -d= -f2-)

build:
	mkdir -p bin
	go build -ldflags="-s -w" -o $(BINARY) ./cmd/server

run:
	go run ./cmd/server

test:
	go test ./... -v -race

test-cover:
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

generate:
	go generate ./...

migrate-up:
	migrate -path $(MIGRATIONS_PATH) -database "$(DB_URL)" up

migrate-down:
	migrate -path $(MIGRATIONS_PATH) -database "$(DB_URL)" down 1

tidy:
	go mod tidy

lint:
	golangci-lint run ./...

docker-up:
	docker compose up --build

docker-down:
	docker compose down -v
