APP_NAME=git_crawler_api
LINTER_VERSION=v1.55.2
MOCKERY_VERSION=v2.30.1
MODULE_NAME := $(shell go list -m)

## -- CI context --

## Run the ci rules
## call [ci-tool, dep, lint, build]
.PHONY: all
all: ci-tool dep lint build

## Install the tooling needs for run the project (ci context):
## - golangci-lint
.PHONY: ci-tool
ci-tool:
# add +"x" in the if statement because if command is not installed VERSION should empty and will cause a syntax error
	@VERSION=$(shell golangci-lint version 2>/dev/null | sed -rn "s/.* (v[0-9]+.[0-9]+.[0-9]+) .*$$/\1/p"); \
	if [ $$VERSION+"x" != ${LINTER_VERSION}+"x" ]; then \
		echo "golangci-lint installation (${LINTER_VERSION})"; \
		go install github.com/golangci/golangci-lint/cmd/golangci-lint@${LINTER_VERSION}; \
		golangci-lint version; \
	fi

## Install tool needs for the development project after calling the ci-tool rule
## - mockery
## - golang.org/x/tools dependencies
.PHONY: tool
tool: ci-tool
	@echo "golang.org/x/tools update ..."
	go install golang.org/x/tools/...@latest
	@VERSION=$(shell mockery --version --quiet 2>/dev/null | sed -rn "s/.*(v[0-9]+.[0-9]+.[0-9]+)$$/\1/p"); \
	if [ $$VERSION+"x" != ${MOCKERY_VERSION}+"x" ]; then \
		echo "mockery installation (${MOCKERY_VERSION})"; \
		go install github.com/vektra/mockery/v2/...@${MOCKERY_VERSION}; \
		mockery --version --quiet; \
	fi

## Install dependencies
.PHONY: dep
dep:
	go mod vendor

## Run linter (golangci-lint) on the full code base
## call ci-tool
.PHONY: lint
lint: ci-tool
	golangci-lint run --skip-dirs mocks --skip-files "(^.+)mock_test.go"

## Run the build
## the default GOOS and GOARCH are used.
## the binaries generated should take the name with <cmd_name>-<GOOS>-<GOARCH>
.PHONY: build
build: 
	$(eval GOOS := $(shell go env GOOS))
	$(eval GOARCH := $(shell go env GOARCH))
	CGO_ENABLED=1 go build -tags static -ldflags "-s -w" \
		-o ${APP_NAME}-${GOOS}-${GOARCH} . 
	@echo "${GREEN}build success" `ls ./${APP_NAME}-${GOOS}-${GOARCH}` "!${NC}"

## -- Dev context --

## Update the generated resources
## should generate or update mocks
.PHONY: update
update:
	GOMODLOCATION=$$PWD go generate ./...

## -- Other commands --

## Cleanup the temporary resources
.PHONY: clean
clean:

## Reset the project to the initial state
## call [clean]
.PHONY: fclean
fclean: clean
	$(eval GOOS := $(shell go env GOOS))
	$(eval GOARCH := $(shell go env GOARCH))
	@rm ${APP_NAME}-${GOOS}-${GOARCH}
#rm -rf vender

