# A Self-Documenting Makefile: http://marmelab.com/blog/2016/02/29/auto-documented-makefile.html

SHELL = /bin/bash
OS = $(shell uname -s)
export PATH := $(abspath bin/):${PATH}

# Project variables
PACKAGE = github.com/banzaicloud/pke

# Build variables
BUILD_DIR ?= $(PWD)/build
BINARY_NAME ?= pke
VERSION ?= $(shell git describe --tags --exact-match 2>/dev/null || git symbolic-ref -q --short HEAD)
COMMIT_HASH ?= $(shell git rev-parse --short HEAD 2>/dev/null)
BUILD_DATE ?= $(shell date +%FT%T%z)
GIT_TREE_STATE ?= $(shell if [[ -z `git status --porcelain 2>/dev/null` ]]; then echo "clean"; else echo "dirty"; fi )
LDFLAGS += -X main.Version=${VERSION} -X main.CommitHash=${COMMIT_HASH} -X main.BuildDate=${BUILD_DATE} -X main.GitTreeState=${GIT_TREE_STATE}
GOPATH ?= `go env GOPATH`

PIPELINE_VERSION = master

# Dependency versions
GOTESTSUM_VERSION = 0.3.4
GOLANGCI_VERSION = 1.26.0
LICENSEI_VERSION = 0.1.0
GORELEASER_VERSION = 0.105.0
OPENAPI_GENERATOR_VERSION = PR1869
TEMPLIFY_VERSION = 7fafacc

GOLANG_VERSION = 1.12

export GOTMPL = ${PWD}/go.tmpl

.PHONY: build
build: pke ## Build project

.PHONY: pke
pke: gogenerate ## Build PKE binary
ifneq (${IGNORE_GOLANG_VERSION_REQ}, 1)
	@printf "${GOLANG_VERSION}\n$$(go version | awk '{sub(/^go/, "", $$3);print $$3}')" | sort -t '.' -k 1,1 -k 2,2 -k 3,3 -g | head -1 | grep -q -E "^${GOLANG_VERSION}$$" || (printf "Required Go version is ${GOLANG_VERSION}\nInstalled: `go version`" && exit 1)
endif
	go build -tags "${GOTAGS}" -ldflags "${LDFLAGS}" -o ${BUILD_DIR}/${BINARY_NAME} ${PACKAGE}/cmd/pke

.PHONY: pke-docs
pke-docs: ## Generate documentation for PKE
	rm -rf cmd/pke/docs/*.md
	cd cmd/pke/docs/ && go run -v generate.go

.PHONY: pke-linux
pke-linux: ## Cross-compile pke for linux
	env GOOS=linux GOARCH=amd64 ${MAKE} pke BINARY_NAME=pke-linux

.PHONY: gogenerate
gogenerate: bin/templify ## Generate go files from template
	GOOS=linux go generate ./cmd/...
	GOOS=darwin go generate ./cmd/...

bin/templify: bin/templify-${TEMPLIFY_VERSION}
	@ln -sf templify-${TEMPLIFY_VERSION} bin/templify
bin/templify-${TEMPLIFY_VERSION}:
	mkdir -p bin
	GOBIN=${PWD}/bin go get -modfile=tools/go.mod github.com/wlbr/templify@${TEMPLIFY_VERSION}
	@mv bin/templify bin/templify-${TEMPLIFY_VERSION}

.PHONY: check
check: test lint ## Run tests and linters

bin/gotestsum: bin/gotestsum-${GOTESTSUM_VERSION}
	@ln -sf gotestsum-${GOTESTSUM_VERSION} bin/gotestsum
bin/gotestsum-${GOTESTSUM_VERSION}:
	@mkdir -p bin
ifeq (${OS}, Darwin)
	curl -L https://github.com/gotestyourself/gotestsum/releases/download/v${GOTESTSUM_VERSION}/gotestsum_${GOTESTSUM_VERSION}_darwin_amd64.tar.gz | tar -zOxf - gotestsum > ./bin/gotestsum-${GOTESTSUM_VERSION} && chmod +x ./bin/gotestsum-${GOTESTSUM_VERSION}
endif
ifeq (${OS}, Linux)
	curl -L https://github.com/gotestyourself/gotestsum/releases/download/v${GOTESTSUM_VERSION}/gotestsum_${GOTESTSUM_VERSION}_linux_amd64.tar.gz | tar -zOxf - gotestsum > ./bin/gotestsum-${GOTESTSUM_VERSION} && chmod +x ./bin/gotestsum-${GOTESTSUM_VERSION}
endif

TEST_PKGS ?= ./...
TEST_REPORT_NAME ?= results.xml
.PHONY: test
test: TEST_REPORT ?= main
test: TEST_FORMAT ?= short
test: SHELL = /bin/bash
test: bin/gotestsum ## Run tests
	@mkdir -p ${BUILD_DIR}/test_results/${TEST_REPORT}
	bin/gotestsum --no-summary=skipped --junitfile ${BUILD_DIR}/test_results/${TEST_REPORT}/${TEST_REPORT_NAME} --format ${TEST_FORMAT} -- $(filter-out -v,${GOARGS}) $(if ${TEST_PKGS},${TEST_PKGS},./...)

bin/golangci-lint: bin/golangci-lint-${GOLANGCI_VERSION}
	@ln -sf golangci-lint-${GOLANGCI_VERSION} bin/golangci-lint
bin/golangci-lint-${GOLANGCI_VERSION}:
	@mkdir -p bin
	curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | bash -s -- -b ./bin/ v${GOLANGCI_VERSION}
	@mv bin/golangci-lint $@

.PHONY: lint
lint: bin/golangci-lint ## Run linter
	bin/golangci-lint run --skip-files '.*\.(json|toml|yaml)\.go'

bin/licensei: bin/licensei-${LICENSEI_VERSION}
	@ln -sf licensei-${LICENSEI_VERSION} bin/licensei
bin/licensei-${LICENSEI_VERSION}:
	@mkdir -p bin
	curl -sfL https://raw.githubusercontent.com/goph/licensei/master/install.sh | bash -s v${LICENSEI_VERSION}
	@mv bin/licensei $@

.PHONY: license-check
license-check: bin/licensei ## Run license check
	bin/licensei check
	./scripts/check-header.sh

.PHONY: license-cache
license-cache: bin/licensei ## Generate license cache
	bin/licensei cache

.PHONY: generate-pipeline-client
generate-pipeline-client: ## Generate client from Pipeline OpenAPI spec
	curl https://raw.githubusercontent.com/banzaicloud/pipeline/${PIPELINE_VERSION}/docs/openapi/pipeline.yaml | sed "s/version: .*/version: ${PIPELINE_VERSION}/" > pipeline-openapi.yaml
	rm -rf .gen/pipeline
	docker run --rm -v ${PWD}:/local banzaicloud/openapi-generator-cli:${OPENAPI_GENERATOR_VERSION} generate \
	--additional-properties packageName=pipeline \
	--additional-properties withGoCodegenComment=true \
	-i /local/pipeline-openapi.yaml \
	-g go \
	-o /local/.gen/pipeline

bin/goreleaser: bin/goreleaser-${GORELEASER_VERSION}
	@ln -sf goreleaser-${GORELEASER_VERSION} bin/goreleaser
bin/goreleaser-${GORELEASER_VERSION}:
	@mkdir -p bin
	curl -sfL https://install.goreleaser.com/github.com/goreleaser/goreleaser.sh | bash -s -- -b ./bin/ v${GORELEASER_VERSION}
	@mv bin/goreleaser $@

.PHONY: release
release: bin/goreleaser # Publish a release
	bin/goreleaser release

# release-%: TAG_PREFIX = v
release-%:
ifneq (${DRY}, 1)
#	@sed -e "s/^## \[Unreleased\]$$/## [Unreleased]\\"$$'\n'"\\"$$'\n'"\\"$$'\n'"## [$*] - $$(date +%Y-%m-%d)/g; s|^\[Unreleased\]: \(.*\/compare\/\)\(.*\)...HEAD$$|[Unreleased]: \1${TAG_PREFIX}$*...HEAD\\"$$'\n'"[$*]: \1\2...${TAG_PREFIX}$*|g" CHANGELOG.md > CHANGELOG.md.new
#	@mv CHANGELOG.md.new CHANGELOG.md

ifeq (${TAG}, 1)
#	git add CHANGELOG.md
#	git commit -m 'Prepare release $*'
	git tag -m 'Release $*' ${TAG_PREFIX}$*
ifeq (${PUSH}, 1)
	git push; git push origin ${TAG_PREFIX}$*
endif
endif
endif

	@echo "Version updated to $*!"
ifneq (${PUSH}, 1)
	@echo
	@echo "Review the changes made by this script then execute the following:"
ifneq (${TAG}, 1)
	@echo
#	@echo "git add CHANGELOG.md && git commit -m 'Prepare release $*' && git tag -m 'Release $*' ${TAG_PREFIX}$*"
	@echo "git tag -m 'Release $*' ${TAG_PREFIX}$*"
	@echo
	@echo "Finally, push the changes:"
endif
	@echo
	@echo "git push; git push origin ${TAG_PREFIX}$*"
endif

.PHONY: patch
patch: ## Release a new patch version
	@${MAKE} release-$(shell (git describe --abbrev=0 --tags 2> /dev/null || echo "0.0.0") | sed 's/^v//' | awk -F'[ .]' '{print $$1"."$$2"."$$3+1}')

.PHONY: minor
minor: ## Release a new minor version
	@${MAKE} release-$(shell (git describe --abbrev=0 --tags 2> /dev/null || echo "0.0.0") | sed 's/^v//' | awk -F'[ .]' '{print $$1"."$$2+1".0"}')

.PHONY: major
major: ## Release a new major version
	@${MAKE} release-$(shell (git describe --abbrev=0 --tags 2> /dev/null || echo "0.0.0") | sed 's/^v//' | awk -F'[ .]' '{print $$1+1".0.0"}')

.PHONY: list
list: ## List all make targets
	@${MAKE} -pRrn : -f $(MAKEFILE_LIST) 2>/dev/null | awk -v RS= -F: '/^# File/,/^# Finished Make data base/ {if ($$1 !~ "^[#.]") {print $$1}}' | egrep -v -e '^[^[:alnum:]]' -e '^$@$$' | sort

.PHONY: help
.DEFAULT_GOAL := help
help:
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

# Variable outputting/exporting rules
var-%: ; @echo $($*)
varexport-%: ; @echo $*=$($*)
