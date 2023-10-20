.DEFAULT_GOAL:=help

## === Tasks ===

## Install binary
install:
	go build -o ${GOPATH}/bin/bruno-http-importer main.go

## Build binary
build:
	mkdir -p bin
	go build -o bin/bruno-http-importer main.go

.PHONY: test

## Run tests
test:
	go test -count=1 ./...

.PHONY: lint
## Run linter
lint:
	golangci-lint run

.PHONY: lint.fix
## Fix lint violations
lint.fix:
	golangci-lint run --fix