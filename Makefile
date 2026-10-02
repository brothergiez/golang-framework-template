.PHONY: help dev build build-api build-cron build-worker build-lambda \
	run-api run-cron run-worker run-lambda test test-unit test-integration \
	lint fmt vet docker-build docker-up docker-down tidy

SERVICE_TYPE ?= api
GO ?= go

help:
	@echo "Targets: build-*, run-*, test, lint, docker-up, docker-down"

tidy:
	$(GO) mod tidy

dev: run-api

build: build-api build-cron build-worker build-lambda

build-api:
	$(GO) build -o bin/api ./cmd/api

build-cron:
	$(GO) build -o bin/cron ./cmd/cron

build-worker:
	$(GO) build -o bin/worker ./cmd/worker

build-lambda:
	$(GO) build -o bin/lambda ./cmd/lambda

run-api:
	SERVICE_TYPE=api DATABASE_DRIVER=memory $(GO) run ./cmd/api

run-cron:
	SERVICE_TYPE=cron DATABASE_DRIVER=memory $(GO) run ./cmd/cron

run-worker:
	SERVICE_TYPE=worker DATABASE_DRIVER=memory $(GO) run ./cmd/worker

run-lambda:
	SERVICE_TYPE=lambda DATABASE_DRIVER=memory $(GO) run ./cmd/lambda

test test-unit:
	$(GO) test ./...

test-integration:
	$(GO) test ./tests/integration/...

lint: vet
	@echo "lint ok"

fmt:
	$(GO) fmt ./...

vet:
	$(GO) vet ./...

docker-build:
	docker build --build-arg SERVICE_TYPE=$(SERVICE_TYPE) -t go-framework:$(SERVICE_TYPE) .

docker-up:
	docker compose up -d api

docker-down:
	docker compose down
