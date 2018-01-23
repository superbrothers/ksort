PROJ=ksort
ORG_PATH=github.com/superbrothers
REPO_PATH=$(ORG_PATH)/$(PROJ)

GIT_COMMIT = $(shell git rev-parse HEAD)
GIT_VERSION = $(shell git describe --tags --abbrev=0 --exact-match 2>/dev/null)
ifeq ($(GIT_VERSION),)
	GIT_VERSION = $(GIT_COMMIT)
endif
GIT_TREE_STATE = $(shell test -n "`git status --porcelain`" && echo "dirty" || echo "clean")
BUILD_DATE = $(shell date -u +'%Y-%m-%dT%H:%M:%SZ')

GO ?= go
OUT_DIR ?= _output
LD_FLAGS :=
LDFLAGS += -X $(REPO_PATH)/version.GitCommit=$(GIT_COMMIT)
LDFLAGS += -X $(REPO_PATH)/version.GitVersion=$(GIT_VERSION)
LDFLAGS += -X $(REPO_PATH)/version.GitTreeState=$(GIT_TREE_STATE)
LDFLAGS += -X $(REPO_PATH)/version.BuildDate=$(BUILD_DATE)

.PHONY: build
build:
		@$(GO) build -o $(OUT_DIR)/$(PROJ) -ldflags '$(LDFLAGS)' .

.PHONY: clean
clean:
		@$(RM) -rf $(OUT_DIR)

HAS_DEP := $(shell command -v dep;)

.PHONY: bootstrap
bootstrap:
ifndef HAS_DEP
	go get -u github.com/golang/dep
endif
