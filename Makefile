.PHONY: all build run test lint db-up db-down clean

# Backend variables
BACKEND_DIR = backend
FRONTEND_DIR = frontend

all: test lint build

build:
	cd $(BACKEND_DIR) && go build -o ../bin/server ./cmd/server
	cd $(FRONTEND_DIR) && npm run build

run:
	cd $(BACKEND_DIR) && go run ./cmd/server

test:
	cd $(BACKEND_DIR) && go test -v -cover ./...

lint:
	cd $(BACKEND_DIR) && golangci-lint run
	cd $(FRONTEND_DIR) && npm run lint

db-up:
	docker compose up -d postgres redis

db-down:
	docker compose stop postgres redis

up:
	docker compose up -d

down:
	docker compose down

clean:
	rm -rf bin/
