# Targets which should always run, regardless of the state of anything else
.PHONY: help install cli test doc

.DEFAULT_GOAL=help

DOCKER_CMD=docker-compose run --rm local
MOUNT=/go/src/github.com/graze/golang-service
CODE=./handlers ./log ./metrics ./nettest

install: ## Install the dependencies
	${DOCKER_CMD} glide install

update: ## Update dependencies
	${DOCKER_CMD} glide update

cli: ## Open a shell to the docker environment
	${DOCKER_CMD} sh

test: ver ?= alpine
test: ## Run the tests
	docker run --rm -it -v $(PWD):${MOUNT} -w ${MOUNT} golang:${ver} go test $$(${DOCKER_CMD} glide nv -x | tr '\n\r' ' ')

doc: ## Build API documentation
	${DOCKER_CMD} godoc github.com/graze/golang-service

lint: ## Run gofmt and goimports in lint mode
	! ${DOCKER_CMD} gofmt -d -s ${CODE} | grep '^'
	! ${DOCKER_CMD} goimports -d ${CODE} | grep '^'

format: ## Run gofmt to format the code
	${DOCKER_CMD} gofmt -s -w ${CODE}
	${DOCKER_CMD} goimports -w ${CODE}

# Build targets
.SILENT: help
help: ## Show this help message
	set -x
	echo "Usage: make [target] ..."
	echo ""
	echo "Available targets:"
	egrep '^(.+)\:\ ##\ (.+)' ${MAKEFILE_LIST} | column -t -c 2 -s ':#' | sort
