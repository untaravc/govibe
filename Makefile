.PHONY: tidy run dev frontend-install frontend-dev frontend-build build
.PHONY: migrate-create migrate-up migrate-down migrate-status migrate-version
.PHONY: seed-all seed-menu

tidy:
	go mod tidy

run:
	go run ./cmd/server

dev:
	air -c air.toml

frontend-install:
	npm install

frontend-dev:
	npm run dev

frontend-build:
	npm run build

migrate-create:
	go run ./cmd/migrate create $(name)

migrate-up:
	go run ./cmd/migrate up

migrate-down:
	go run ./cmd/migrate down

migrate-status:
	go run ./cmd/migrate status

migrate-version:
	go run ./cmd/migrate version

seed-all:
	go run ./cmd/seed all

seed-menu:
	go run ./cmd/seed menu

build:
	mkdir -p bin
	go build -o bin/server ./cmd/server
