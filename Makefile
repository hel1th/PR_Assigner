.PHONY: build run docker-up docker-down test-up test-down test-wait test clean deps

lint: 
	golangci-lint run ./...

build:
	go build -o bin/pr_assigner ./cmd/api

run:
	go run ./cmd/api/main.go

docker-up:
	docker-compose up -d --build

docker-down:
	docker-compose down

docker-logs:
	docker-compose logs -f

test-up:
	docker-compose -f docker-compose.test.yml up -d --build

test-down:
	docker-compose -f docker-compose.test.yml down -v

test-wait:
	@echo "Waiting for test API..."
	@timeout 30 bash -c 'until curl -sf http://localhost:8081/health; do sleep 1; done' && echo "✓ Ready"

test: test-down test-up test-wait
	@echo "Running E2E tests..."
	@DATABASE_URL=postgres://test_user:test_password@localhost:5434/pr_assigner_test?sslmode=disable \
	API_URL=http://localhost:8081 \
	go test ./test/e2e/ -v -count=1
	@make test-down

test-quick: test-wait
	@DATABASE_URL=postgres://test_user:test_password@localhost:5434/pr_assigner_test?sslmode=disable \
	API_URL=http://localhost:8081 \
	go test ./test/e2e/ -v -count=1

deps:
	go mod download
	go mod tidy

clean:
	rm -rf bin/
	go clean -cache
	docker-compose down -v
	docker-compose -f docker-compose.test.yml down -v