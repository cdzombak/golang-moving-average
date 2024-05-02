SHELL:=/usr/bin/env bash

default: help
.PHONY: help
help: ## Print help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: lint
lint: ## Lint all .go files
	go vet ./...
	golangci-lint run *.go

.PHONY: test
test: ## Run tests
	go test -v -race ./...
