.PHONY: build
build:
	go build -o bin/mediaexporter -v ./cmd/main.go

.PHONY: run
run:
	go run ./cmd/main.go

.PHONY: test
test:
	go test -v -timeout 30s ./...

.DEFAULT_GOAL := build