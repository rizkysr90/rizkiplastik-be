.PHONY: vendor
vendor: 
	go mod tidy && go mod vendor
.PHONY: server/start
server/start:
	docker-compose -f docker-compose.yml up -d  --remove-orphans

.PHONY: server/stop
server/stop:
	docker-compose -f docker-compose.yml down

.PHONY: server/restart
server/restart:
	docker-compose -f docker-compose.yml restart api db-migration
################
# BUILD BINARY
################

GOOS ?= linux
GOARCH ?= amd64
CGO_ENABLED ?= 0

.PHONY: build/restapi
build/restapi:
	CGO_ENABLED=$(CGO_ENABLED) GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o build/restapi cmd/main.go

.PHONY: test
test:
	go test ./internal/service/... -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html

GOLANGCI_VERSION=1.55.2
GOLANGCI_CHECK := $(shell golangci-lint -v 2> /dev/null)

.PHONY: lint
lint:
# if golangci-lint failed on MacOS Ventura, try: brew install diffutils
ifndef GOLANGCI_CHECK
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v$(GOLANGCI_VERSION)
endif
	golangci-lint run -c .golangci.yml ./...