SHELL = /bin/bash

# Project variables
PACKAGE = github.com/banzaicloud/pke

# Build variables
BUILD_DIR ?= $(PWD)/build
VERSION ?= $(shell git describe --tags --exact-match 2>/dev/null || git symbolic-ref -q --short HEAD)
COMMIT_HASH ?= $(shell git rev-parse --short HEAD 2>/dev/null)
BUILD_DATE ?= $(shell date +%FT%T%z)
GIT_TREE_STATE ?= $(shell if [[ -z `git status --porcelain 2>/dev/null` ]]; then echo "clean"; else echo "dirty"; fi )
LDFLAGS += -X main.Version=${VERSION} -X main.CommitHash=${COMMIT_HASH} -X main.BuildDate=${BUILD_DATE} -X main.GitTreeState=${GIT_TREE_STATE}

.PHONY: pke
pke: GOARGS += -tags "${GOTAGS}" -ldflags "${LDFLAGS}" -o ${BUILD_DIR}/pke
pke: ## Build PKE binary
	cd cmd/pke/ && go build ${GOARGS} ${PACKAGE}/cmd/pke

.PHONY: pke-docs
pke-docs: ## Generate documentation for PKE
	rm -rf cmd/pke/docs/*.md
	cd cmd/pke/docs/ && go run -v generate.go

.PHONY: test
test: export CGO_ENABLED = 1
test:
	set -o pipefail; go list ./... | xargs -n1 go test ${GOARGS} -v -parallel 1 2>&1 | tee test.txt


.PHONY: help
.DEFAULT_GOAL := help
help:
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
