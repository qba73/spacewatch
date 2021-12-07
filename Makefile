.PHONY: help check cover test tidy vet run

ROOT := $(PWD)
GO_HTML_COV := ./coverage.html
GO_TEST_OUTFILE := ./c.out
GO_DOCKER_IMAGE := golang:1.17
GO_DOCKER_CONTAINER := meteo-container


define PRINT_HELP_PYSCRIPT
import re, sys

for line in sys.stdin:
	match = re.match(r'^([a-zA-Z_-]+):.*?## (.*)$$', line)
	if match:
		target, help = match.groups()
		print("%-20s %s" % (target, help))
endef
export PRINT_HELP_PYSCRIPT

default: help

help:
	@python -c "$$PRINT_HELP_PYSCRIPT" < $(MAKEFILE_LIST)

vet: ## Run go vet and shadow
	go vet ./...
	shadow ./...

check: ## Run static check analyzer
	staticcheck ./...

cover: ## Run unit tests and generate test coverage report
	go test -shuffle=on -race -v ./... -count=1 -cover -covermode=atomic -coverprofile=coverage.out
	go tool cover -html coverage.out
	staticcheck ./...

test: vet ## Run unit tests locally
	go test -shuffle=on -race -v ./...
	staticcheck ./...

# MODULES
tidy: ## Run go mod tidy and vendor
	go mod tidy
	go mod vendor


# DEVELOPMENT
run: ## Run service locally
	go run cmd/spacewatch-api/main.go
