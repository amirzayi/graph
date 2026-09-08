.PHONY: test coverage

.DEFAULT_GOAL := deps

deps:
	go mod tidy

up:
	docker compose up

down:
	docker compose down

test:
	go mod tidy
	go test -v -count=1 ./...

run:
	go mod tidy
	go run .

coverage:
	go mod tidy
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out
