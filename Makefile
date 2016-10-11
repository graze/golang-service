# Targets which should always run, regardless of the state of anything else
.PHONY: help

.DEFAULT_GOAL=install

DOCKER_CMD=docker-compose run --rm local
ver=alpine
PATH=/go/src/github.com/graze/golang-service
TEST_CMD=docker run --rm -it -v $(PWD):${PATH} -w ${PATH} golang:${ver}

install: ## Install the dependencies
	${DOCKER_CMD} glide install

cli: ## Open a shell to the docker environment
	${DOCKER_CMD} sh

test: ## Run all tests
	${TEST_CMD} go test ./nettest ./logging ./handlers

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
