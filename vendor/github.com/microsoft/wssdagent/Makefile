# Copyright (c) Microsoft Corporation.
# Licensed under the MIT license.
GOCMD=go
GOBUILD=$(GOCMD) build -v -mod=vendor
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

endif

format:
	gofmt -s -w cmd/ common/ rpc/ test/ rpc/ pkg/ services/

agent:
	GO111MODULE=on GOARCH=amd64 GOOS=windows $(GOBUILD) -ldflags "-X main.version=$(TAG) -X main.commit=$(COMMIT) -X main.date=$(BUILD_DATE)" -o ${OUT} github.com/microsoft/wssdagent/cmd/wssdagent
