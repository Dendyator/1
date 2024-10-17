BIN := "./bin/calendar"
DOCKER_IMG := "calendar:develop"
GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

build:
	go build -v -o $(BIN) -ldflags "$(LDFLAGS)" ./cmd/calendar

run-db:
	docker-compose -f deployments/docker-compose.yaml up -d db

# Функция ожидания, чтобы убедиться, что база данных работает
wait-db:
	@echo "Waiting for database to be ready..."
	@until docker exec postgres_db pg_isready -U user; do \
		echo "Waiting for database..."; \
		sleep 2; \
	done

run: run-db build-img
	@if [ "$(docker ps -q -f name=calendar_service)" ]; then \
		docker stop calendar_service; \
		docker rm calendar_service; \
	fi; \
	docker run --network calendar_network -p 8080:8080 --name calendar_service --rm $(DOCKER_IMG)

build-img:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_IMG) \
		-f deployments/build/Dockerfile .

run-img: run-db build-img
	docker run -p 8080:8080 --name calendar_service --rm $(DOCKER_IMG)

version: build
	$(BIN) version

test:
	go test -race -cover ./internal/...

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.61.0

lint: install-lint-deps
	golangci-lint run ./...
	go fmt ./...
	go vet ./...

migrate:
	goose -dir ./migrations postgres "your_database_connection_string" up

.PHONY: build run run-db build-img run-img version test lint wait-db


