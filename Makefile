.PHONY: build
build:
	go build -o omniserve ./cmd/omniserve

.PHONY: install
install: build
	mv omniserve /usr/local/bin/

.PHONY: clean
clean:
	rm -f omniserve