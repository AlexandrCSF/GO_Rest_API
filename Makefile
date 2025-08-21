.PHONY: build
build:
	go build -v ./cmd/apiserver

.PHONY: test
test:
	go test -v -race -timeout=30s ./...

.PHONY: run
run:
	go run ./cmd/apiserver

.PHONY: emulator
emulator:
	go run ./scripts/kafka_emulator.go

.PHONY: migrate
migrate:
	psql -h localhost -U postgres -d restapi_dev -f migrations/20241222140000_create_orders.up.sql

.PHONY: rollback
rollback:
	psql -h localhost -U postgres -d restapi_dev -f migrations/20241222140000_create_orders.down.sql

.DEFAULT_GOAL := build