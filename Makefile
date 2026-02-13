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

.PHONY: build all test test-unit test-bdd lint fmt-check vet check-prek ci clean release-snapshot audit help

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

## check-prek: Run pre-commit checks when prek.toml exists
check-prek:
	@if [ -f prek.toml ]; then \
		command -v prek >/dev/null 2>&1 || { echo "prek.toml exists but 'prek' is not installed" >&2; exit 1; }; \
		prek validate-config; \
		prek run --all-files; \
	fi

## fmt-check: Ensure gofmt is clean
fmt-check:
	@unformatted=$$(gofmt -l $$(find . -type f -name '*.go' -not -path './bin/*' -not -path './dist/*')); \
	if [ -n "$$unformatted" ]; then \
		echo "Unformatted files:"; \
		echo "$$unformatted"; \
		exit 1; \
	fi

## vet: Run go vet
vet:
	go vet ./...

## test: Run all tests (unit + BDD), with pre-commit checks first
test: check-prek test-unit test-bdd

## test-unit: Run unit tests
test-unit:
	go test -v -race ./internal/... ./cmd/...

## test-bdd: Build binary and run BDD tests
test-bdd: build
	CIO_BINARY=$(CURDIR)/$(BINDIR)/$(BINARY) go test -v ./test/...

## lint: Run golangci-lint
lint:
	golangci-lint run ./...

## ci: Full local/CI quality gate
ci: check-prek fmt-check vet test lint

## release-snapshot: Local GoReleaser dry-run (no publish)
release-snapshot:
	goreleaser build --snapshot --clean

## audit: Audit workflows and release naming contract
audit:
	scripts/audit-workflows.sh .
	scripts/audit-release-naming.sh .

## clean: Remove build artifacts
clean:
	rm -rf $(BINDIR)/ dist/ cio

## help: Show this help
help:
	@grep -E '^## ' $(MAKEFILE_LIST) | sed 's/## //' | column -t -s ':'
