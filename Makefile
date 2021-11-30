.PHONY:
.SILENT:
.DEFAULT_GOAL := run

build:
	go mod download && CGO_ENABLED=0 GOOS=linux go build -o ./.bin/app ./cmd/app/main.go

run:
	go run ./cmd/app

swag:
	swag init -g internal/app/app.go

lint:
	golangci-lint run
