# Copyright (c) Microsoft Corporation.
# Licensed under the MIT license.
GOCMD=go
GOBUILD=$(GOCMD) build -i -v -mod=vendor
GOHOSTOS=$(strip $(shell $(GOCMD) env get GOHOSTOS))

TAG ?= $(shell git describe --tags)
COMMIT ?= $(shell git describe --always)
BUILD_DATE ?= $(shell date -u +%m/%d/%Y)

OUT=bin/wssdagent.exe

PKG := 

all: agent

clean:
	rm -rf ${OUT} 
.PHONY: vendor
vendor:
	GO111MODULE=on go mod vendor
ifeq ($(GOHOSTOS),linux)
	git clone --branch vm https://github.com/madhanrm/hcsshim.git
	rm -rf vendor/github.com/Microsoft/hcsshim
	mv hcsshim vendor/github.com/Microsoft/.

	git clone https://github.com/census-instrumentation/opencensus-go vendor/go.opencensus.io

	mkdir -p vendor/github.com/hashicorp
	git clone https://github.com/hashicorp/golang-lru vendor/github.com/hashicorp/golang-lru

	mkdir vendor/github.com/golang/groupcache
	git clone https://github.com/golang/groupcache vendor/github.com/golang/groupcache

endif

format:
	gofmt -s -w cmd/ common/ rpc/ test/ rpc/ pkg/ services/

agent:
	GO111MODULE=on GOARCH=amd64 GOOS=windows $(GOBUILD) -ldflags "-X main.version=$(TAG) -X main.commit=$(COMMIT) -X main.date=$(BUILD_DATE)" -o ${OUT} github.com/microsoft/wssdagent/cmd/wssdagent
