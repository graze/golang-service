# Targets which should always run, regardless of the state of anything else
.PHONY: help

.DEFAULT_GOAL=help

DOCKER_CMD=docker-compose run --rm golang
VERSION=`git describe --tags`
LDFLAGS=-ldflags "-X github.com/graze/logging/version.version=${VERSION}"

install: ## Install dependencies
	mkdir -p bin
	${DOCKER_CMD} go get \
		github.com/DataDog/datadog-go/statsd \
		github.com/gorilla/handlers \
		github.com/stretchr/testify/assert

cli:
	${DOCKER_CMD} sh

test: ## Run all tests
	${DOCKER_CMD} go test ./logging ./nettest

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
