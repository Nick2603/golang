.PHONY: all build run swagger tidy docker-build docker-up docker-down install-swag

all: tidy swagger build

install-swag:
	go install github.com/swaggo/swag/cmd/swag@latest

tidy:
	go mod tidy

swagger: install-swag
	swag init -g cmd/server/main.go --output docs

build: swagger
	go build -o bin/server ./cmd/server/main.go

run: swagger
	go run ./cmd/server/main.go

docker-build:
	docker-compose build

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

docker-logs:
	docker-compose logs -f app
