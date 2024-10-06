APPLICATION	:= kycelisd
NAMESPACE		?= nicklasfrahm
REGISTRY		?= docker.io
VERSION			?= $(shell git describe --tags --always --dirty)
IMAGE				:= $(REGISTRY)/$(NAMESPACE)/$(APPLICATION):$(VERSION)
GOFLAGS			:= -ldflags "-X main.version=$(VERSION) -s -w"
UPXFLAGS		?= ""

define HELP_HEADER
Usage:	make <target>

Targets:
endef

export HELP_HEADER

help: ## List all targets.
	@echo "$$HELP_HEADER"
	@grep -E '^[a-zA-Z0-9%_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'

.env: ## Create a seed .env file.
	touch .env

.PHONY: run
run: .env ## Run the kycelisd server locally.
	LOG_FORMAT=console go run $(GOFLAGS) cmd/$(APPLICATION)/main.go

.PHONY: build
build: ## Build the kycelisd server.
	CGO_ENABLED=0 go build $(GOFLAGS) -o bin/$(APPLICATION) cmd/$(APPLICATION)/main.go
ifneq ($(UPXFLAGS),"")
	upx $(UPXFLAGS) bin/$(APPLICATION)
endif

.PHONY: test
test: ## Run tests.
	go test -cover -v ./...

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

.PHONY: container-build
container-build: ## Build the container.
	docker build \
		--build-arg VERSION=$(VERSION) \
		-f build/package/Containerfile \
		-t $(IMAGE) .

.PHONY: container-push
container-push: ## Push the container.
	docker push $(IMAGE)

.PHONY: container-export
container-export: ## Export the container.
	@mkdir -p dist
	docker save $(IMAGE) -o dist/$(APPLICATION).tar

.PHONY: container-import
container-import: ## Import the container.
	docker load -i dist/$(APPLICATION).tar

.PHONY: infra-plan
infra-plan: ## Plan the infrastructure.
	tofu -chdir=deploy/tofu init
	tofu -chdir=deploy/tofu plan -var-file=local.tfvars

.PHONY: infra-apply
infra-apply: ## Apply the infrastructure.
	tofu -chdir=deploy/tofu init
	tofu -chdir=deploy/tofu apply -var-file=local.tfvars -auto-approve
