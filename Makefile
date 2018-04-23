PWD := $(shell pwd)
APP := frog
PKG := github.com/ckeyer/$(APP)

GO := go
HASH := $(shell which shasum || which sha1sum)

OS := $(shell go env GOOS)
ARCH := $(shell go env GOARCH)
VERSION := $(shell cat VERSION.txt)
GIT_COMMIT := $(shell git rev-parse --short HEAD)
GIT_BRANCH := $(shell git rev-parse --abbrev-ref HEAD)
BUILD_AT := $(shell date "+%Y-%m-%dT%H:%M:%SZ%z")
PACKAGE_NAME := $(APP)$(VERSION).$(OS)-$(ARCH)

LD_FLAGS := -X github.com/ckeyer/commons/version.version=$(VERSION) \
 -X github.com/ckeyer/commons/version.gitCommit=$(GIT_COMMIT) \
 -X github.com/ckeyer/commons/version.buildAt=$(BUILD_AT) -w

GO_IMAGE := ckeyer/go:1.10

env:
	$(GO) env

gorun:
	$(GO) run -ldflags="$(LD_FLAGS)" main.go

build-in-docker:
	docker run --rm -it \
	 -e CGO_ENABLED=0 \
	 -v `pwd`:/go/src/${PKG} \
	 -v `pwd`/bundles:/go/bin/ \
	 -w /go/src/${PKG} \
	 ${GO_IMAGE} make build

build: env
	$(GO) build -v -ldflags="$(LD_FLAGS)" -o ${GOPATH}/bin/$(APP) main.go
	$(HASH) ${GOPATH}/bin/$(APP)

test:
	$(GO) test $$(go list ./... |grep -v "vendor")

test-in-docker:
	docker run --rm \
	 -v ${PWD}:/go/src/${PKG} \
	 -w /go/src/${PKG} \
	 ${GO_IMAGE} \
	 go test -ldflags="$(LD_FLAGS)" $$(go list ./... |grep -v "vendor")

clean:
	rm -rf bundles/*

dev:
	docker run --rm -it \
	 --name $(APP)-dev \
	 -p 8080:8080 \
	 -v $(PWD)/..:/opt/gopath/src/$(PKG)/.. \
	 -w /opt/gopath/src/$(PKG) \
	 $(GO_IMAGE) sh
