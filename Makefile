# Copyright (c) Microsoft Corporation.
# Licensed under the MIT license.
GOCMD=go
GOBUILD=$(GOCMD) build -v #-mod=vendor
GOHOSTOS=$(strip $(shell $(GOCMD) env get GOHOSTOS))
GOPATH_BIN := $(shell go env GOPATH)/bin
GOTEST=GOOS=$(GOHOSTOS) $(GOCMD) test -v -coverprofile=coverage.out -covermode count -timeout 60m0s
TESTDIRECTORIES= ./services/compute/virtualmachine/internal ./services/security/keyvault/key/internal

TAG ?= $(shell git describe --tags)
COMMIT ?= $(shell git describe --always)
BUILD_DATE ?= $(shell date -u +%m/%d/%Y)
LDFLAGS="-X main.version=$(TAG) -X main.commit=$(COMMIT) -X main.date=$(BUILD_DATE)"

# Private repo workaround
export GOPRIVATE=github.com/microsoft
# Active module mode, as we use go modules to manage dependencies
export GO111MODULE=on

LBCLIENTOUT=bin/lbclient.exe

all: format  lbclient build unittest

nofmt: 

clean:
	rm -rf 	${LBCLIENTOUT} 
lbclient:
	GOARCH=amd64 GOOS=windows $(GOBUILD) -ldflags $(LDFLAGS) -o ${LBCLIENTOUT} github.com/microsoft/wssd-sdk-for-go/cmd/lbclient
format:
	gofmt -s -w cmd/ pkg/ services/ tools/

.PHONY: vendor
vendor:
	go mod tidy
build:
	GOARCH=amd64 go build -v ./...

unittest:
	$(GOTEST) $(TESTDIRECTORIES)

golangci-lint:
	$(GOCMD) install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	$(GOPATH_BIN)/golangci-lint run --config .golangci.yml

	