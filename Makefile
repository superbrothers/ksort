PROJ := ksort
ORG_PATH := github.com/superbrothers
REPO_PATH := $(ORG_PATH)/$(PROJ)

GIT_COMMIT := $(shell git rev-parse HEAD)
GIT_VERSION := $(shell git describe --tags --abbrev=0 --exact-match 2>/dev/null)
ifeq ($(GIT_VERSION),)
GIT_VERSION := $(GIT_COMMIT)
endif
GIT_TREE_STATE := $(shell test -n "`git status --porcelain`" && echo "dirty" || echo "clean")
BUILD_DATE := $(shell date -u +'%Y-%m-%dT%H:%M:%SZ')

GO_VERSION ?= 1.10
GOOS ?= $(shell uname | tr A-Z a-z)
GOARCH ?= amd64
GOCACHE ?= $(shell pwd)/.go-build
GO_WORKDIR := /go/src/$(REPO_PATH)
GO ?= docker run --rm -e GOOS -e GOARCH -e CGO_ENABLED=0 -w $(GO_WORKDIR) -v $(shell pwd):$(GO_WORKDIR) -v $(GOCACHE):/root/.cache/go-build golang:$(GO_VERSION) go
OUT_DIR ?= _output
LD_FLAGS :=
LDFLAGS += -X $(REPO_PATH).GitCommit=$(GIT_COMMIT)
LDFLAGS += -X $(REPO_PATH).GitVersion=$(GIT_VERSION)
LDFLAGS += -X $(REPO_PATH).GitTreeState=$(GIT_TREE_STATE)
LDFLAGS += -X $(REPO_PATH).BuildDate=$(BUILD_DATE)

.PHONY: build
build:
		@$(GO) build -o $(OUT_DIR)/$(PROJ) -a -installsuffix cgo -ldflags '$(LDFLAGS)' ./cmd/ksort

.PHONY: test
test:
		@$(GO) test -v ./...

.PHONY: clean
clean:
		@$(RM) -rf $(OUT_DIR)

HAS_DEP := $(shell command -v dep;)

.PHONY: bootstrap
bootstrap:
ifndef HAS_DEP
	go get -u github.com/golang/dep
endif
