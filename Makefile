SHELL = /bin/bash

BUF_CMD := $(shell pwd)/bin/buf
BUF_VERSION := latest
GOOS := Linux
PE_SUFFIX :=

ifeq ($(OS), Windows_NT)
	GOOS := Windows
	PE_SUFFIX := .exe
else ifeq ($(shell uname -s), Darwin)
	GOOS := Darwin
endif

export PATH := $(shell pwd)/bin:$(PATH)

SPEC_FILE ?= spec/openapi.yaml
BUNDLE_PATH ?= dist/openapi.yaml

openapi-dependencies: ## Install openapi dependencies
	cd api/openapi && npm install

openapi-generate: ## Generate OpenAPI
	@./scripts/utils/openapi-http.sh client src/go/internal/genopenapi/server server

openapi-bundle: openapi-dependencies ## Generate OpenAPI bundle
	cd api/openapi && npm run bundle -- ${SPEC_FILE} -o ${BUNDLE_PATH}

openapi-preview: ## Preview OpenAPI
	cd api/openapi &&  SPEC_FILE=$(SPEC_FILE) npm run preview

buf-dependencies: ## Install dependencies
	mkdir -p bin
	curl -sSLo bin/buf$(PE_SUFFIX) https://github.com/bufbuild/buf/releases/$(BUF_VERSION)/download/buf-$(GOOS)-x86_64$(PE_SUFFIX)
	curl -sSLo bin/protoc-gen-buf-breaking$(PE_SUFFIX) https://github.com/bufbuild/buf/releases/$(BUF_VERSION)/download/protoc-gen-buf-breaking-$(GOOS)-x86_64$(PE_SUFFIX)
	curl -sSLo bin/protoc-gen-buf-lint$(PE_SUFFIX) https://github.com/bufbuild/buf/releases/$(BUF_VERSION)/download/protoc-gen-buf-lint-$(GOOS)-x86_64$(PE_SUFFIX)
	chmod +x bin/*buf*

buf-generate: ## Generate proto
	cd api/proto && $(BUF_CMD) generate

.DEFAULT_GOAL := help
.PHONY: help
help: ## Show this help screen.
	@echo 'Usage: make <OPTIONS> ... <TARGETS>'
	@echo ''
	@echo 'Available targets are:'
	@echo ''
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z0-9_-]+:.*?##/ { printf "  \033[36m%-25s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)
