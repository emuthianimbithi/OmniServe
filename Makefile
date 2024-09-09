VERSION ?= $(shell git describe --tags --abbrev=0)

.PHONY: build
build:
	go build -ldflags "-X main.Version=$(VERSION)" -o omniserve ./cmd/omniserve

.PHONY: install
install: build
	mv omniserve /usr/local/bin/

.PHONY: clean
clean:
	rm -f omniserve
