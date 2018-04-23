PWD := $(shell pwd)
APP := frog
PKG := github.com/ckeyer/$(APP)

GO := go
HASH := $(shell which shasum || which sha1sum)

OS := $(shell go env GOOS)
ARCH := $(shell go env GOARCH)
VERSION := $(shell cat VERSION)
GIT_COMMIT := $(shell git rev-parse --short HEAD)
GIT_BRANCH := $(shell git rev-parse --abbrev-ref HEAD)
BUILD_AT := $(shell date "+%Y-%m-%dT%H:%M:%SZ%z")
PACKAGE_NAME := $(APP)$(VERSION).$(OS)-$(ARCH)

COMMONS_PKG := $(PKG)/vendor/github.com/ckeyer/commons
LD_FLAGS := -X ${COMMONS_PKG}/version.version=$(VERSION) \
 -X ${COMMONS_PKG}/version.gitCommit=$(GIT_COMMIT) \
 -X ${COMMONS_PKG}/version.buildAt=$(BUILD_AT)

GO_IMAGE := ckeyer/go:1.10

env:
	$(GO) env

gorun:
	$(GO) run -ldflags="$(LD_FLAGS)" main.go

build-in-docker:
	docker run --rm -it \
	 -e CGO_ENABLED=0 \
	 -v `pwd`:/go/src/${PKG} \
	 -w /go/src/${PKG} \
	 ${GO_IMAGE} make build

build:
	$(GO) build -v -ldflags="$(LD_FLAGS)" -o bundles/$(APP) main.go
	$(HASH) bundles/$(APP)

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
	 -v $(PWD):/go/src/$(PKG) \
	 -w /go/src/$(PKG) \
	 $(GO_IMAGE) sh
