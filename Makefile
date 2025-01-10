
BINARY_NAME=myapp


SRC_DIR=./cmd
MIGRATE_DIR=./migrations


GO=go
GOFMT=gofmt
GOFLAGS=-v
LDFLAGS=-s -w
GOLINT=golangci-lint 

migrate:
	$(GO) run $(MIGRATE_DIR)

lint:
	$(GOLINT) run ./...


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

all: test build
