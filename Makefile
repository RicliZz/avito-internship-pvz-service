.PHONY: build run migrate-create migrate-up migrate-down

include .env
export

build:
	@go build -o ./cmd/bin ./cmd/main.go
run:build
	@go run ./cmd/main.go
migrate:
	@migrate create -ext sql -dir db/migrations -seq ${name}
migrate-up:
	@migrate -database ${POSTGRESQL_URL} -path db/migrations up

migrate-down:
	@migrate -database ${POSTGRESQL_URL} -path db/migrations down
