BINARY_NAME=myapp

SRC_DIR=./cmd
MIGRATE_DIR=./migrations
DOCKER_COMPOSE_DIR=./deployments/docker-compose.yml

GO=go
GOFMT=gofmt
GOFLAGS=-v
LDFLAGS=-s -w
GOLINT=golangci-lint
DOCKERCOMPOSE=docker compose

migrate:
	$(GO) run $(MIGRATE_DIR)

lint:
	$(GOLINT) run ./...

container-build:
	$(DOCKERCOMPOSE) -f $(DOCKER_COMPOSE_DIR) up --build;

container-run:
	$(DOCKERCOMPOSE) -f $(DOCKER_COMPOSE_DIR) up

run:
	$(GO) run $(SRC_DIR)

deps:
	$(GO) mod tidy

format:
	$(GOFMT) -w .

runf: run format
precommit: lint format
all: test build
