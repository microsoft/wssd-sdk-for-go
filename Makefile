# Copyright (c) Microsoft Corporation.
# Licensed under the MIT license.
GOCMD=go
GOBUILD=$(GOCMD) build -v -mod=vendor
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
	GO111MODULE=on GOARCH=amd64 GOOS=windows $(GOBUILD) -ldflags "-X main.version=$(TAG) -X main.commit=$(COMMIT) -X main.date=$(BUILD_DATE)" -o ${OUT} github.com/microsoft/wssd-sdk-for-go/cmd/wssdctl

.PHONY: vendor
vendor:
	GO111MODULE=on go get github.com/microsoft/wssdagent 
	GO111MODULE=on go mod vendor
	rm -rf vendor/github.com/Microsoft/hcsshim
	git clone --branch vm https://github.com/madhanrm/hcsshim.git vendor/github.com/Microsoft/hcsshim
	git clone https://github.com/census-instrumentation/opencensus-go vendor/go.opencensus.io
	mkdir -p vendor/github.com/hashicorp
	git clone https://github.com/hashicorp/golang-lru vendor/github.com/hashicorp/golang-lru
