.PHONY: run build test clean deps docker-up docker-down

run:
	go run cmd/api/main.go

build:
	go build -o bin/api cmd/api/main.go

test:
	go test -v ./tests/integration/...

test-coverage:
	go test -v -coverprofile=coverage.out ./tests/integration/...
	go tool cover -html=coverage.out

clean:
	rm -rf bin/ docs/ coverage.out

deps:
	go mod download
	go mod tidy

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

docker-logs:
	docker-compose logs -f
