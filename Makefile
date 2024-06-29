# Include: minimum
-include .makefiles/minimum/Makefile
.makefiles/minimum/Makefile:
	@git clone https://github.com/tmknom/makefiles.git .makefiles >/dev/null 2>&1

# Variables: Go
VERSION ?= 0.0.1
ROOT_DIR ?= $(shell \git rev-parse --show-toplevel)
NAME = $(shell \basename $(ROOT_DIR))
COMMIT = $(shell \git rev-parse HEAD)
DATE = $(shell \date +"%Y-%m-%d")
LDFLAGS ?= "-X main.name=$(NAME) -X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.date=$(DATE)"

# Targets: Go
.PHONY: all
all: mod build lint test run ## all

.PHONY: mod
mod: ## manage modules
	go mod tidy
	go mod verify

.PHONY: deps
deps:
	go mod download

.PHONY: build
build: deps ## build executable binary
	go build -ldflags=$(LDFLAGS) -o bin/tfmod ./cmd/tfmod

.PHONY: install
install: deps ## install
	go install -ldflags=$(LDFLAGS) ./cmd/tfmod

.PHONY: run
run: build ## run command
	@bin/tfmod

.PHONY: test
test: deps ## test all
	go test ./...

.PHONY: lint
lint: goimports vet ## lint go

.PHONY: vet
vet: ## static analysis by vet
	go vet ./...

.PHONY: goimports
goimports: ## update import lines
	goimports -w .

.PHONY: install-tools
install-tools: ## install tools for development
	go install golang.org/x/tools/cmd/goimports@latest

# Targets: GitHub Actions
.PHONY: lint-gha
lint-gha: lint/workflow lint/yaml ## Lint workflow files and YAML files

.PHONY: fmt-gha
fmt-gha: fmt/yaml ## Format YAML files
