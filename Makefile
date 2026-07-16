PREFIX ?= /usr/local

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
