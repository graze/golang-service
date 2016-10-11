# Targets which should always run, regardless of the state of anything else
.PHONY: help install cli test doc

.DEFAULT_GOAL=help

DOCKER_CMD=docker-compose run --rm local
MOUNT=/go/src/github.com/graze/golang-service

install: ## Install the dependencies
	${DOCKER_CMD} glide install

cli: ## Open a shell to the docker environment
	${DOCKER_CMD} sh

test: ver ?= alpine
test: ## Run the tests
	docker run --rm -it -v $(PWD):${MOUNT} -w ${MOUNT} golang:${ver} go test ./nettest ./logging ./handlers

doc: ## Build API documentation
	${DOCKER_CMD} godoc github.com/graze/golang-service

# Build targets
.SILENT: help
help: ## Show this help message
	set -x
	echo "Usage: make [target] ..."
	echo ""
	echo "Available targets:"
	egrep '^(.+)\:\ ##\ (.+)' ${MAKEFILE_LIST} | column -t -c 2 -s ':#' | sort
