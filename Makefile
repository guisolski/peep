PREFIX ?= $(shell test -d /usr/local/bin && echo /usr/local || echo /opt/homebrew)
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
LDFLAGS := -X main.version=$(VERSION)

.DEFAULT_GOAL := help

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-12s\033[0m %s\n", $$1, $$2}'

build: ## Build the peep binary (./peep)
	go build -ldflags "$(LDFLAGS)" -o peep .

install: build ## Build and install to $(PREFIX)/bin
	install -m 0755 peep $(PREFIX)/bin/peep

uninstall: ## Remove the installed binary from $(PREFIX)/bin
	rm -f $(PREFIX)/bin/peep

vet: ## Run go vet
	go vet ./...

fmt: ## Check formatting with gofmt
	gofmt -l -s .

lint: ## Run golangci-lint
	golangci-lint run ./...

test: ## Run tests with the race detector
	go test ./... -race

bench: ## Run benchmarks
	go test ./... -run=^$$ -bench=. -benchmem

docs-serve: ## Serve the docs locally with mkdocs
	mkdocs serve

clean: ## Remove the built binary
	rm -f peep

.PHONY: help build install uninstall vet fmt lint test bench docs-serve clean
