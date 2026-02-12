BINARY  := cio
BINDIR  := bin
MODULE  := github.com/leechael/cio
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
COMMIT  := $(shell git rev-parse --short HEAD 2>/dev/null || echo none)
DATE    := $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
LDFLAGS := -s -w \
  -X $(MODULE)/cmd.version=$(VERSION) \
  -X $(MODULE)/cmd.commit=$(COMMIT) \
  -X $(MODULE)/cmd.date=$(DATE)

PLATFORMS := \
  linux/amd64 \
  linux/arm64 \
  darwin/amd64 \
  darwin/arm64 \
  windows/amd64

.PHONY: build test test-unit test-bdd lint clean all release-snapshot help

## build: Build for current platform
build:
	CGO_ENABLED=0 go build -ldflags "$(LDFLAGS)" -o $(BINDIR)/$(BINARY) .

## all: Cross-compile for all platforms
all: $(PLATFORMS)

$(PLATFORMS):
	$(eval GOOS := $(word 1,$(subst /, ,$@)))
	$(eval GOARCH := $(word 2,$(subst /, ,$@)))
	$(eval EXT := $(if $(filter windows,$(GOOS)),.exe,))
	CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) \
		go build -ldflags "$(LDFLAGS)" -o dist/$(BINARY)-$(GOOS)-$(GOARCH)$(EXT) .

## test: Run all tests (unit + BDD)
test: test-unit test-bdd

## test-unit: Run unit tests
test-unit:
	go test -v -race ./internal/... ./cmd/...

## test-bdd: Build binary and run BDD tests
test-bdd: build
	CIO_BINARY=$(CURDIR)/$(BINDIR)/$(BINARY) go test -v ./test/...

## lint: Run golangci-lint
lint:
	golangci-lint run ./...

## release-snapshot: Local GoReleaser dry-run (no publish)
release-snapshot:
	goreleaser build --snapshot --clean

## clean: Remove build artifacts
clean:
	rm -rf $(BINDIR)/ dist/

## help: Show this help
help:
	@grep -E '^## ' $(MAKEFILE_LIST) | sed 's/## //' | column -t -s ':'
