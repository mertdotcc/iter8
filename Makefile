BINDIR      := $(CURDIR)/bin
INSTALL_PATH ?= /usr/local/bin
DIST_DIRS   := find * -type d -exec
TARGETS     ?= darwin/amd64 linux/amd64 linux/386 windows/amd64
BINNAME     ?= iter8

GOBIN         = $(shell go env GOBIN)
ifeq ($(GOBIN),)
GOBIN         = $(shell go env GOPATH)/bin
endif
GOX           = $(GOBIN)/gox

# go options
TAGS       :=
LDFLAGS    := -w -s
GOFLAGS    :=

# Rebuild the binary if any of these files change
SRC := $(shell find . -type f -name '*.(go|proto|tpl)' -print) go.mod go.sum

# Required for globs to work correctly
SHELL      = /usr/bin/env bash

GIT_COMMIT = $(shell git rev-parse HEAD)
GIT_TAG    = $(shell git describe --tags --dirty)

ifdef VERSION
	BINARY_VERSION = $(VERSION)
endif
BINARY_VERSION ?= ${GIT_TAG}

# Only set Version if GIT_TAG or VERSION is set
ifneq ($(BINARY_VERSION),)
	LDFLAGS += -X github.com/iter8-tools/iter8/base.Version=${BINARY_VERSION}
endif


LDFLAGS += -X github.com/iter8-tools/iter8/cmd.gitCommit=${GIT_COMMIT}

.PHONY: all
all: build

# ------------------------------------------------------------------------------
#  build

.PHONY: build
build: $(BINDIR)/$(BINNAME)

$(BINDIR)/$(BINNAME): $(SRC)
	GO111MODULE=on go build $(GOFLAGS) -trimpath -tags '$(TAGS)' -ldflags '$(LDFLAGS)' -o '$(BINDIR)'/$(BINNAME) ./

# ------------------------------------------------------------------------------
#  install

.PHONY: install
install: build
	@install "$(BINDIR)/$(BINNAME)" "$(INSTALL_PATH)/$(BINNAME)"

# ------------------------------------------------------------------------------
#  dependencies

$(GOX):
	go install github.com/mitchellh/gox@latest

# ------------------------------------------------------------------------------
#  release

.PHONY: build-cross
build-cross: LDFLAGS += -extldflags "-static"
build-cross: $(GOX)
	GOFLAGS="-trimpath" GO111MODULE=on CGO_ENABLED=0 $(GOX) -parallel=2 -output="_dist/{{.OS}}-{{.Arch}}/$(BINNAME)" -osarch='$(TARGETS)' $(GOFLAGS) -tags '$(TAGS)' -ldflags '$(LDFLAGS)' ./

.PHONY: dist
dist: build-cross
	( \
		cd _dist && \
		$(DIST_DIRS) cp ../LICENSE {} \; && \
		$(DIST_DIRS) cp ../README.md {} \; && \
		$(DIST_DIRS) tar -zcf iter8-{}.tar.gz {} \; \
	)

.PHONY: clean
clean:
	@rm -rf '$(BINDIR)' ./_dist

# ------------------------------------------------------------------------------
#  test

.PHONY: fmt
fmt: ## Run go fmt against code.
	go fmt ./...

.PHONY: vet
vet: ## Run go vet against code
	go vet ./...

.PHONY: staticcheck
staticcheck:
	staticcheck ./...

.PHONY: golangci-lint
golangci-lint:
	golangci-lint run ./...

.PHONY: lint
lint:	vet staticcheck golangci-lint

.PHONY: test
test: fmt vet ## Run tests.
	go test -v ./... -coverprofile=coverage.out

.PHONY: coverage
coverage: test
	@echo "test coverage: $(shell go tool cover -func coverage.out | grep total | awk '{print substr($$3, 1, length($$3)-1)}')"

.PHONY: htmlcov
htmlcov: coverage
	go tool cover -html=coverage.out
