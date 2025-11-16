build:
	go build -o bin/pr_assigner ./cmd/api

run:
	go run ./cmd/api/main.go

docker-up:
	docker-compose up -d --build

docker-down:
	docker-compose down

deps:
	go mod download
	go mod tidy

clean:
	rm -rf bin/
	go clean -cache

.PHONY: build run docker-up docker-down clean deps