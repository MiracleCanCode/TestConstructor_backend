
BINARY_NAME=myapp


SRC_DIR=./cmd
MIGRATE_DIR=./migrations
DOCKER_COMPOSE_DIR=./deployments/docker-compose.yml


GO=go
GOFMT=gofmt
GOFLAGS=-v
LDFLAGS=-s -w
GOLINT=golangci-lint 
DOCKERCOMPOSE=docker-compose

migrate:
	$(GO) run $(MIGRATE_DIR)

lint:
	$(GOLINT) run ./...

docker-compose:
	$(DOCKERCOMPOSE) -f $(DOCKER_COMPOSE_DIR) up --build

build:
	$(GO) build $(GOFLAGS) -o $(BINARY_NAME) $(SRC_DIR)

run:
	$(GO) run $(SRC_DIR)

test:
	$(GO) test $(GOFLAGS) ./...

benchmark:
	$(GO) test -benchmem -bench .

deps:
	$(GO) mod tidy

format:
	$(GOFMT) -w .

runf: run format
precommit: lint format
all: test build
