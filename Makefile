SHELL = /bin/bash

# Configurable Variables
BUF_VERSION := latest
GOOS := $(shell uname -s)
PE_SUFFIX := $(if $(findstring Windows,$(GOOS)),.exe,)
BIN_DIR := $(shell pwd)/bin
SPEC_FILE ?= spec/openapi.yaml
BUNDLE_PATH ?= dist/openapi.yaml

# Export PATH
export PATH := $(BIN_DIR):$(PATH)

# Targets
.PHONY: all openapi-dependencies openapi-generate openapi-bundle openapi-preview buf-dependencies buf-generate help

all: help ## Default target

## OpenAPI Tasks
openapi-dependencies: ## Install openapi dependencies
	cd api/openapi && npm install

openapi-generate: ## Generate OpenAPI
	@./scripts/utils/openapi-http.sh src/go/pkg/api/codegen codegen

openapi-bundle: openapi-dependencies ## Generate OpenAPI bundle
	cd api/openapi && npm run bundle -- ${SPEC_FILE} -o ${BUNDLE_PATH}

openapi-preview: ## Preview OpenAPI
	cd api/openapi &&  SPEC_FILE=$(SPEC_FILE) npm run preview

## Buf/Protobuf Tasks
buf-dependencies: ## Install Buf dependencies
	mkdir -p $(BIN_DIR)
	curl -sSLo $(BIN_DIR)/buf$(PE_SUFFIX) https://github.com/bufbuild/buf/releases/$(BUF_VERSION)/download/buf-$(GOOS)-x86_64$(PE_SUFFIX)
	curl -sSLo $(BIN_DIR)/protoc-gen-buf-breaking$(PE_SUFFIX) https://github.com/bufbuild/buf/releases/$(BUF_VERSION)/download/protoc-gen-buf-breaking-$(GOOS)-x86_64$(PE_SUFFIX)
	curl -sSLo $(BIN_DIR)/protoc-gen-buf-lint$(PE_SUFFIX) https://github.com/bufbuild/buf/releases/$(BUF_VERSION)/download/protoc-gen-buf-lint-$(GOOS)-x86_64$(PE_SUFFIX)
	chmod +x $(BIN_DIR)/*buf*

buf-generate: ## Generate Protobuf code
	cd api/proto && $(BIN_DIR)/buf generate

## Help Target
help: ## Show this help screen
	@echo 'Usage: make <target>'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z0-9_-]+:.*?##/ { printf "  \033[36m%-25s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)