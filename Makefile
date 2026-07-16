PREFIX ?= $(shell test -d /usr/local/bin && echo /usr/local || echo /opt/homebrew)

build:
	go build -o peep .

install: build
	install -m 0755 peep $(PREFIX)/bin/peep

uninstall:
	rm -f $(PREFIX)/bin/peep

test:
	go test ./...

clean:
	rm -f peep

.PHONY: build install uninstall test clean
