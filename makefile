# ====================================================================================
# VARIABLES
# ====================================================================================
BINARY_NAME := martini
BUILD_DIR := bin
MAIN_PACKAGE := .

VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)

# System install (Unix layout). Use: sudo make install
PREFIX ?= /usr/local
BINDIR := $(PREFIX)/bin

GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)
CGO_ENABLED ?= 0

LDFLAGS := -s -w
BUILD_FLAGS := -trimpath

GOCMD := go
GOBUILD := $(GOCMD) build
GOCLEAN := $(GOCMD) clean
GOTEST := $(GOCMD) test
GOINSTALL := $(GOCMD) install

GOBIN := $(shell go env GOBIN)
ifeq ($(GOBIN),)
GOBIN := $(shell go env GOPATH)/bin
endif

BUILD_LDFLAGS := $(LDFLAGS)

# ====================================================================================
# DEFAULT TARGET
# ====================================================================================
.DEFAULT_GOAL := help

# ====================================================================================
# PHONY TARGETS
# ====================================================================================
.PHONY: help all tidy fmt vet test check build run clean \
	install install-user uninstall uninstall-user

help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Development:"
	@echo "  tidy          Tidy go.mod and go.sum"
	@echo "  fmt           Format Go source"
	@echo "  vet           Run go vet"
	@echo "  test          Run tests with coverage"
	@echo "  check         fmt, vet, and test"
	@echo "  build         Build $(BINARY_NAME) to $(BUILD_DIR)/"
	@echo "  run           Run via go run"
	@echo "  clean         Remove build artifacts"
	@echo ""
	@echo "Install:"
	@echo "  install-user  Install to user bin via go install ($(GOBIN)/)"
	@echo "  install       Install to $(BINDIR)/ (often: sudo make install)"
	@echo "  uninstall-user  Remove from $(GOBIN)/"
	@echo "  uninstall     Remove from $(BINDIR)/ (often: sudo make uninstall)"

all: check build

tidy:
	$(GOCMD) mod tidy

fmt:
	$(GOCMD) fmt ./...

vet: fmt
	$(GOCMD) vet ./...

test:
	$(GOTEST) -race -coverprofile=coverage.out ./...

check: vet test

build:
	@mkdir -p $(BUILD_DIR)
	CGO_ENABLED=$(CGO_ENABLED) GOOS=$(GOOS) GOARCH=$(GOARCH) \
		$(GOBUILD) $(BUILD_FLAGS) -ldflags="$(BUILD_LDFLAGS)" \
		-o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PACKAGE)

run:
	$(GOCMD) run $(MAIN_PACKAGE)

clean:
	$(GOCLEAN)
	rm -rf $(BUILD_DIR)
	rm -f coverage.out

# User install: respects GOBIN and GOPATH/bin (same idea as ./build.sh)
install-user:
	@echo "Installing $(BINARY_NAME) via go install to user PATH..."
	CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) \
		$(GOINSTALL) $(BUILD_FLAGS) -ldflags="$(BUILD_LDFLAGS)" $(MAIN_PACKAGE)
	@echo "Installed to $(GOBIN)/$(BINARY_NAME) (ensure $(GOBIN) is on your PATH)"

# System install: copy built binary into PREFIX (DESTDIR for packaging roots)
install: build
	@echo "Installing $(BINARY_NAME) to $(DESTDIR)$(BINDIR)..."
	@mkdir -p $(DESTDIR)$(BINDIR)
	@install -m 755 $(BUILD_DIR)/$(BINARY_NAME) $(DESTDIR)$(BINDIR)/$(BINARY_NAME)
	@echo "Done: $(DESTDIR)$(BINDIR)/$(BINARY_NAME)"

uninstall-user:
	@echo "Removing $(GOBIN)/$(BINARY_NAME)..."
	@rm -f $(GOBIN)/$(BINARY_NAME)
	@echo "Done."

uninstall:
	@echo "Removing $(BINDIR)/$(BINARY_NAME)..."
	@rm -f $(BINDIR)/$(BINARY_NAME)
	@echo "Done."
