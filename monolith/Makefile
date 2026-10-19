ifneq (,$(wildcard ./.env))
    include .env
    export
endif

POSTGRES_USER ?= spearszx
POSTGRES_PASSWORD ?= secret_pass1234
POSTGRES_HOST ?= localhost
POSTGRES_DB ?= go_wallet
POSTGRES_PORT ?= 5438

MIGRATIONS_PATH ?= ./db/migrations

DB_DSN = postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=disable

.PHONY: build run test docker-up docker-stop migrate-up migrate-down migrate-create .env cover

cover:
	go test -short -count=1 -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out
	rm coverage.out

.env:
	@cp .env.example .env
	@echo "Created .env from .env.example"


build:
	@go build -o ./bin/app ./cmd/main.go

#=================================================================================================================
#
#	Запуск из корневой директории (для синхронизации с файлом .env.example)
#
#=================================================================================================================
run: build
	@./bin/app


test:
	@go test -v -count=1 ./...


docker-up:
	docker-compose up -d


docker-down:
	docker-compose down -v


migrate-up:
	goose -dir $(MIGRATIONS_PATH) postgres $(DB_DSN) up


migrate-down:
	goose -dir $(MIGRATIONS_PATH) postgres $(DB_DSN) down


migrate-create:
	@test -n "$(name)" || (echo "Usage: make migrate-create name=<migration_name>" && exit 1)
	@echo "$(name)" | grep -Eq '^[a-zA-Z0-9_-]+$$' || (echo "Invalid name: use only letters, digits, _, -" && exit 1)
	goose -dir $(MIGRATIONS_PATH) create $(name) sql
