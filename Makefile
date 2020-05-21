.PHONY: build
build:
	go build -o bin/imageresize -v ./cmd/imageresize

.PHONY: test
test:
	go test -v -race -timeout 30s ./...

.PHONY: run
run:
	./bin/imageresize

.DEFAULT_GOAL := build