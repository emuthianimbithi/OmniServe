export

VERSION ?= $(shell git describe --tags --abbrev=0 2>/dev/null || echo "dev")

.PHONY: build
build:
	go build -ldflags "-X main.Version=$(VERSION)" -o omniserve ./cmd/omniserve

.PHONY: build-all
build-all: build-linux build-macos build-windows

.PHONY: build-linux
build-linux:
	GOOS=linux GOARCH=amd64 go build -ldflags "-X main.Version=$(VERSION)" -o omniserve-linux-amd64 ./cmd/omniserve

.PHONY: build-macos
build-macos:
	GOOS=darwin GOARCH=amd64 go build -ldflags "-X main.Version=$(VERSION)" -o omniserve-macos-amd64 ./cmd/omniserve

.PHONY: build-windows
build-windows:
	GOOS=windows GOARCH=amd64 go build -ldflags "-X main.Version=$(VERSION)" -o omniserve-windows-amd64.exe ./cmd/omniserve

.PHONY: install
install: build
	mv omniserve /usr/local/bin/

.PHONY: clean
clean:
	rm -f omniserve omniserve-linux-amd64 omniserve-macos-amd64 omniserve-windows-amd64.exe

proto:
	 protoc --go_out=pkg/pb --go_opt=paths=source_relative \
                    --go-grpc_out=pkg/pb --go-grpc_opt=paths=source_relative \
                    -I=pkg/pb \
                    pkg/pb/**/*.proto
