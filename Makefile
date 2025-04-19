.PHONY: build run migrate-create migrate-up migrate-down

include .env
export

build:
	@go build -o ./cmd/mainApp/bin ./cmd/mainApp/main.go
run:build
	./cmd/mainApp/bin
migrate:
	@migrate create -ext sql -dir db/migrations -seq ${name}
migrate-up:
	@migrate -database ${POSTGRESQL_URL} -path db/migrations up

migrate-down:
	@migrate -database ${POSTGRESQL_URL} -path db/migrations down ${count}
