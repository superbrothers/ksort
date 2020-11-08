GO ?= GO111MODULE=on GOPROXY=https://gocenter.io go
DIST_DIR := dist

.PHONY: build
build:
	$(GO) build -o $(DIST_DIR)/ksort cmd/ksort/main.go

TOOLS_DIR := hack/tools
TOOLS_BIN_DIR := $(TOOLS_DIR)/bin
GORELEASER_BIN := bin/goreleaser
GORELEASER := $(TOOLS_DIR)/$(GORELEASER_BIN)
GOLANGCI_LINT_BIN := bin/golangci-lint
GOLANGCI_LINT := $(TOOLS_DIR)/$(GOLANGCI_LINT_BIN)

$(GORELEASER): $(TOOLS_DIR)/go.mod
	cd $(TOOLS_DIR) && $(GO) build -o $(GORELEASER_BIN) github.com/goreleaser/goreleaser

$(GOLANGCI_LINT): $(TOOLS_DIR)/go.mod
	cd $(TOOLS_DIR) && $(GO) build -o $(GOLANGCI_LINT_BIN) github.com/golangci/golangci-lint/cmd/golangci-lint


.PHONY: build-cross
build-cross: $(GORELEASER)
	$(GORELEASER) build --snapshot --rm-dist

.PHONY: vet
vet:
	$(GO) vet ./...

.PHONY: fmt
fmt:
	$(GO) fmt ./...

.PHONY: lint
lint: $(GOLANGCI_LINT) vet fmt
	$(GOLANGCI_LINT) run

.PHONY: test
test:
	$(GO) test -v ./...

.PHONY: dist
dist: $(GORELEASER)
	$(GORELEASER) release --rm-dist --skip-publish --snapshot

.PHONY: release
release: $(GORELEASER)
	$(GORELEASER) release --rm-dist

.PHONY: clean
clean:
	rm -rf $(DIST_DIR) $(TOOLS_BIN_DIR)
