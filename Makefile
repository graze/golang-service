# Targets which should always run, regardless of the state of anything else
.PHONY: help

.DEFAULT_GOAL=help

DOCKER_CMD=docker-compose run --rm golang
DOCKER_CMD_FROM_ROOT=docker-compose run -w /go golang
VERSION=`git describe --tags`
LDFLAGS=-ldflags "-X github.com/graze/golang-service/version.version=${VERSION}"

build:
	docker-compose build

cli:
	${DOCKER_CMD} sh

test: ## Run all tests
	${DOCKER_CMD} go test ./...

doc: ## Build API documentation
	${DOCKER_CMD} godoc github.com/graze/golang-service

# Build targets
.SILENT: help
help: ## Show this help message
	set -x
	echo "Usage: make [target] ..."
	echo ""
	echo "Available targets:"
	fgrep --no-filename "##" $(MAKEFILE_LIST) | fgrep --invert-match $$'\t' | sed -e 's/: ## / - /'
