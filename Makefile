define HELP_HEADER
Usage:	make <target>

Targets:
endef

export HELP_HEADER

help: ## List all targets.
	@echo "$$HELP_HEADER"
	@grep -E '^[a-zA-Z0-9%_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-10s\033[0m %s\n", $$1, $$2}'

.PHONY: test
test: ## Run tests.
	go test -v ./...

LINTER := $(GOPATH)/bin/golangci-lint
$(LINTER):
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

.PHONY: lint
lint: $(LINTER) ## Run linter.
	$< run

FORMATTER := $(GOPATH)/bin/gofumpt
$(FORMATTER):
	go install mvdan.cc/gofumpt@latest

.PHONY: format
format: $(FORMATTER) ## Format the code.
	$< -l -w .
