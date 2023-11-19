.PHONY: build
build:
	go build -o bin/mediaexporter -v ./cmd/main.go

.PHONY: run
run:
	go run ./cmd/main.go

.PHONY: test
test:
	CGO_ENABLED=1 go test -race -v -timeout 30s ./...

MOCKS_PREFIX=mocks
.PHONY: mocks
mocks: internal/adapters/media/yi4kplus/telnet/telnetclient.go
	@echo "Generating mocks..."
	@rm -rf $(MOCKS_DESTINATION)
	@for file in $^; do mockgen -source=$$file -destination=$(MOCKS_PREFIX)/$$file; done

#TODO: написать команду генерации моков, на основе этого пример:
#mockgen -source=telnetclient.go -destination=$(MOCKS_PREFIX)/telnetclient.go -package=$(MOCKS_PREFIX)
#mockgen -source=telnetclient.go -destination=mocks/telnetclient.go -package=mocks

.DEFAULT_GOAL := build