.PHONY: deps up down test run coverage bench profiling

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

bench:
	go test -bench=. -run=^$ -benchmem ./handler


profiling:
	go test -bench=. -run=^$ -benchmem -cpuprofile=cpu.prof -memprofile=mem.prof ./handler
	go tool pprof -http=:8083 cpu.prof &
	go tool pprof -http=:8084 mem.prof
