.PHONY: help build run test clean docker-up docker-down migrate-up migrate-down seed deps

help:
	@echo "Laundry Management System - Makefile Commands"
	@echo ""
	@echo "build          - Build the API binary"
	@echo "run            - Run the API locally (requires PostgreSQL)"
	@echo "test           - Run unit tests"
	@echo "clean          - Remove build artifacts"
	@echo "docker-up      - Start PostgreSQL via docker-compose"
	@echo "docker-down    - Stop PostgreSQL"
	@echo "migrate-up     - Run all pending migrations"
	@echo "migrate-down   - Roll back all migrations"
	@echo "seed           - Seed the database with initial data"
	@echo "deps           - Download and tidy Go dependencies"

build:
	go build -o bin/laundry-api ./cmd/api/main.go

run:
	go run ./cmd/api/main.go

test:
	go test -v ./...

clean:
	rm -rf bin/
	go clean

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

migrate-up:
	go run ./cmd/api/main.go --migrate-only

migrate-down:
	@echo "Applying down migrations..."
	@for f in $$(ls migrations/*.down.sql | sort -r); do \
		echo "Rolling back: $$f"; \
		psql "$$DATABASE_URL" -f "$$f"; \
	done

seed:
	go run ./scripts/seed.go

deps:
	go mod download
	go mod tidy
