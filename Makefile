# Copyright (c) Microsoft Corporation.
# Licensed under the MIT license.
GOCMD=go
GOBUILD=$(GOCMD) build -v #-mod=vendor
GOHOSTOS=$(strip $(shell $(GOCMD) env get GOHOSTOS))

TAG ?= $(shell git describe --tags)
COMMIT ?= $(shell git describe --always)
BUILD_DATE ?= $(shell date -u +%m/%d/%Y)
LDFLAGS="-X main.version=$(TAG) -X main.commit=$(COMMIT) -X main.date=$(BUILD_DATE)"

# Private repo workaround
export GOPRIVATE = github.com/microsoft
# Active module mode, as we use go modules to manage dependencies
export GO111MODULE=on

LBCLIENTOUT=bin/lbclient.exe

PKG := 

all: format  lbclient

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
