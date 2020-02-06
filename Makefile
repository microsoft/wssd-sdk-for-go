# Copyright (c) Microsoft Corporation.
# Licensed under the MIT license.
GOCMD=go
GOBUILD=$(GOCMD) build -v -mod=vendor
GOHOSTOS=$(strip $(shell $(GOCMD) env get GOHOSTOS))

TAG ?= $(shell git describe --tags)
COMMIT ?= $(shell git describe --always)
BUILD_DATE ?= $(shell date -u +%m/%d/%Y)

OUTEXE=bin/wssdctl.exe
OUT=bin/wssdctl

PKG := 

all: format ctl ctlexe

clean:
	rm -rf ${OUT} ${OUTEXE}
ctlexe:
	GO111MODULE=on GOARCH=amd64 GOOS=windows $(GOBUILD) -ldflags "-X main.version=$(TAG) -X main.commit=$(COMMIT) -X main.date=$(BUILD_DATE)" -o ${OUTEXE} github.com/microsoft/wssd-sdk-for-go/cmd/wssdctl
ctl:
	GO111MODULE=on GOARCH=amd64 $(GOBUILD) -ldflags "-X main.version=$(TAG) -X main.commit=$(COMMIT) -X main.date=$(BUILD_DATE)" -o ${OUT} github.com/microsoft/wssd-sdk-for-go/cmd/wssdctl
format:
	gofmt -s -w cmd/ common/ pkg/ services/ test/ tools/

.PHONY: vendor
vendor:
	#GO111MODULE=on GOPRIVATE="github.com/microsoft" go get github.com/microsoft/wssdagent 
	GO111MODULE=on go mod vendor
	GO111MODULE=on go mod tidy
