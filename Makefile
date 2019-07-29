# Copyright (c) Microsoft Corporation.
# Licensed under the MIT license.
GOCMD=go
GOBUILD=$(GOCMD) build -i -v #-mod=vendor
GOHOSTOS=$(strip $(shell $(GOCMD) env get GOHOSTOS))

TAG ?= $(shell git describe --tags)
COMMIT ?= $(shell git describe --always)
BUILD_DATE ?= $(shell date -u +%m/%d/%Y)

OUT=bin/wssdctl.exe

PKG := 

all: ctl

clean:
	rm -rf ${OUT} 
ctl:
	GOARCH=amd64 GOOS=windows $(GOBUILD) -ldflags "-X main.version=$(TAG) -X main.commit=$(COMMIT) -X main.date=$(BUILD_DATE)" -o ${OUT} github.com/microsoft/wssd-sdk-for-go/cmd/wssdctl

.PHONY: vendor

vendor:
	go get all ./...
