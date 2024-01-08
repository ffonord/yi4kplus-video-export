.PHONY: build
build:
	go build -o bin/mediaexporter -v ./cmd/main.go

.PHONY: run
run:
	go run ./cmd/main.go

.PHONY: test
test: mocks
	CGO_ENABLED=1 go test -race -v -timeout 30s ./...

MOCKS_PREFIX=mocks
.PHONY: mocks
mocks: internal/adapters/media/yi4kplus/amba/ambaclient.go \
	internal/adapters/media/yi4kplus/ftp/ftpclient.go \
	internal/adapters/media/yi4kplus/telnet/telnetclient.go
	@echo "Generating mocks..."
	@for filepath in $^; \
		do \
		  	dst=$$(dirname $$filepath)/$(MOCKS_PREFIX)/$$(basename $$filepath); \
			mockgen -package=$(MOCKS_PREFIX) -source=$$filepath -destination=$${dst}; \
		done

.DEFAULT_GOAL := build