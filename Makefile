PREFIX ?= $(shell test -d /usr/local/bin && echo /usr/local || echo /opt/homebrew)
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
LDFLAGS := -X main.version=$(VERSION)

build:
	go build -ldflags "$(LDFLAGS)" -o peep .

install: build
	install -m 0755 peep $(PREFIX)/bin/peep

uninstall:
	rm -f $(PREFIX)/bin/peep

vet:
	go vet ./...

fmt:
	gofmt -l -s .

lint:
	golangci-lint run ./...

test:
	go test ./... -race

clean:
	rm -f peep

.PHONY: build install uninstall vet fmt lint test clean
