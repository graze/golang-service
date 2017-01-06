# Targets which should always run, regardless of the state of anything else
.PHONY: help install cli test doc

.DEFAULT_GOAL=help

DOCKER_CMD=docker-compose run --rm tools
MOUNT=/go/src/github.com/graze/golang-service
CODE=./handlers ./handlers/auth ./handlers/recovery ./log ./metrics ./nettest ./validate

install: ## Install the dependencies
	rm -rf vendor
	${DOCKER_CMD} glide install

update: ## Update dependencies
	rm -rf vendor
	${DOCKER_CMD} glide update

cli: ## Open a shell to the docker environment
	${DOCKER_CMD} sh

test: ver ?= alpine
test: ## Run the tests
	docker run --rm -it -v $(PWD):${MOUNT} -w ${MOUNT} golang:${ver} go test ${CODE}

doc: ## Build API documentation
	${DOCKER_CMD} godoc github.com/graze/golang-service

lint: ## Run gofmt and goimports in lint mode
	${DOCKER_CMD} golint -set_exit_status ./handlers/...
	${DOCKER_CMD} golint -set_exit_status ./log/...
	${DOCKER_CMD} golint -set_exit_status ./metrics/...
	${DOCKER_CMD} golint -set_exit_status ./nettest/...
	${DOCKER_CMD} golint -set_exit_status ./validate/...
	${DOCKER_CMD} golint -set_exit_status ./
	${DOCKER_CMD} go tool vet ./handlers
	${DOCKER_CMD} go tool vet ./log
	${DOCKER_CMD} go tool vet ./metrics
	${DOCKER_CMD} go tool vet ./nettest
	${DOCKER_CMD} go tool vet ./validate

format: ## Run gofmt to format the code
	${DOCKER_CMD} gofmt -s -w ${CODE}
	${DOCKER_CMD} goimports -w ${CODE}

clean: ## Clean docker and git info
	docker-compose stop
	docker-compose rm -f
	docker-compose down --remove-orphans || echo "Cleaned"
	git clean -d -f -f
	rm -rf .glide

# Build targets
.SILENT: help
help: ## Show this help message
	set -x
	echo "Usage: make [target] ..."
	echo ""
	echo "Available targets:"
	egrep '^(.+)\:\ ##\ (.+)' ${MAKEFILE_LIST} | column -t -c 2 -s ':#' | sort
