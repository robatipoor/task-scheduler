APP_NAMES := master worker
VERSION := $(shell git describe --tags --abbrev=0)
BUILD_DIR := build
SRC_DIR := cmd
PKG_DIR := pkg
INTERNAL_DIR := internal
TEST_DIR := test
GO := go
GOFMT := gofmt
GOTEST := go test
GOBUILD := $(GO) build
GOCLEAN := $(GO) clean


.PHONY: all
all: build

.PHONY: build
build: clean $(APP_NAMES)

$(APP_NAMES):
	@echo "Building the application $@..."
	$(GOBUILD) -o $(BUILD_DIR)/$@ $(SRC_DIR)/$@/main.go

.PHONY: test
test:
	@echo "Running tests..."
	$(GOTEST) ./... -v

.PHONY: test-coverage
test-coverage:
	@echo "Running tests with coverage..."
	$(GOTEST) -cover ./... -v

.PHONY: fmt
fmt:
	@echo "Formatting the code..."
	$(GOFMT) -w .

.PHONY: clean
clean:
	@echo "Cleaning up..."
	$(GOCLEAN)
	rm -rf $(BUILD_DIR)

.PHONY: deps
deps:
	@echo "Installing dependencies..."
	$(GO) mod tidy

.PHONY: run-master
run-master: build
	@echo "Running application master..."
	./$(BUILD_DIR)/master

.PHONY: run-worker
run-worker: build
	@echo "Running application worker..."
	./$(BUILD_DIR)/worker

.PHONY: help
help:
	@echo "Makefile for $(APP_NAMES)"
	@echo ""
	@echo "Usage:"
	@echo "  make [target]"
	@echo ""
	@echo "Targets:"
	@echo "  all              Build the application (default)"
	@echo "  build            Build the application"
	@echo "  test             Run tests"
	@echo "  test-coverage    Run tests with coverage"
	@echo "  fmt              Format the code"
	@echo "  clean            Clean build artifacts"
	@echo "  deps             Install dependencies"
	@echo "  run-master       Run the master node"
	@echo "  run-worker       Run the worker node"
	@echo "  help             Show this help message"
