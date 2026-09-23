# ====================================================================================
# VARIABLES
# ====================================================================================
BINARY_NAME := martini
BUILD_DIR := bin
MAIN_PACKAGE := .

GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)
CGO_ENABLED ?= 0

LDFLAGS := -s -w
BUILD_FLAGS := -trimpath

GOBIN := $(shell go env GOBIN)
ifeq ($(GOBIN),)
GOBIN := $(shell go env GOPATH)/bin
endif

# ====================================================================================
# DEFAULT TARGET
# ====================================================================================
.DEFAULT_GOAL := help

# ====================================================================================
# PHONY TARGETS
# ====================================================================================
.PHONY: help tidy fmt vet test check build install run clean

help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@echo "  tidy      Tidy go.mod and go.sum"
	@echo "  fmt       Format Go source"
	@echo "  vet       Run go vet"
	@echo "  test      Run tests with coverage"
	@echo "  check     fmt, vet, and test"
	@echo "  build     Build $(BINARY_NAME) to $(BUILD_DIR)/"
	@echo "  install   Install static $(BINARY_NAME) to $(GOBIN)/"
	@echo "  run       Run via go run (development)"
	@echo "  clean     Remove build artifacts"

tidy:
	go mod tidy

fmt:
	go fmt ./...

vet: fmt
	go vet ./...

test:
	go test -race -coverprofile=coverage.out ./...

check: vet test

build:
	@mkdir -p $(BUILD_DIR)
	CGO_ENABLED=$(CGO_ENABLED) GOOS=$(GOOS) GOARCH=$(GOARCH) \
		go build $(BUILD_FLAGS) -ldflags="$(LDFLAGS)" \
		-o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PACKAGE)

install:
	@mkdir -p $(GOBIN)
	CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) \
		go build $(BUILD_FLAGS) -ldflags="$(LDFLAGS)" \
		-o $(GOBIN)/$(BINARY_NAME) $(MAIN_PACKAGE)
	@echo "Installed $(GOBIN)/$(BINARY_NAME)"

run:
	go run $(MAIN_PACKAGE)

clean:
	go clean
	rm -rf $(BUILD_DIR)
	rm -f coverage.out
