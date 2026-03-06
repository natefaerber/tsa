BINARY  := tsa
MODULE  := github.com/natefaerber/tsa
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
COMMIT  := $(shell git rev-parse --short HEAD 2>/dev/null || echo unknown)
LDFLAGS := -s -w -X $(MODULE)/cmd.Version=$(VERSION) -X $(MODULE)/cmd.Commit=$(COMMIT)

.PHONY: all build test lint fmt vet clean install

all: test build

build:
	go build -ldflags "$(LDFLAGS)" -o $(BINARY) .

test:
	go test ./... -v

lint: vet
	@command -v golangci-lint >/dev/null 2>&1 || { echo "golangci-lint not installed"; exit 1; }
	golangci-lint run ./...

fmt:
	gofmt -l -w .

vet:
	go vet ./...

clean:
	rm -f $(BINARY)
	go clean

install: build
	go install -ldflags "$(LDFLAGS)" .
