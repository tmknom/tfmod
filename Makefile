# Include: minimum
-include .makefiles/minimum/Makefile
.makefiles/minimum/Makefile:
	@git clone https://github.com/tmknom/makefiles.git .makefiles >/dev/null 2>&1

# Variables: Go
NAME = $(shell \basename $(ROOT_DIR))
VERSION = $(shell \git tag --sort=-v:refname | head -1)
OWNER = $(shell \gh config get -h github.com user)
COMMIT = $(shell \git rev-parse HEAD)
DATE = $(shell \date +"%Y-%m-%d")
URL = https://github.com/$(OWNER)/$(NAME)/releases/tag/$(VERSION)
LDFLAGS ?= "-X main.name=$(NAME) -X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.date=$(DATE) -X main.url=$(URL)"

# Targets: Go
.PHONY: all
all: mod build test-all run ## all

.PHONY: mod
mod: ## manage modules
	go mod tidy
	go mod verify

.PHONY: deps
deps:
	go mod download

.PHONY: build
build: deps ## build executable binary
	go build -ldflags=$(LDFLAGS) -o bin/$(NAME) ./cmd/$(NAME)

.PHONY: install
install: deps ## install
	go install -ldflags=$(LDFLAGS) ./cmd/$(NAME)

.PHONY: run
run: build ## run command
	bin/tfmod get --debug --base-dir "testdata/terraform/env"
	@echo
	bin/tfmod dependencies --debug --base-dir "testdata/terraform" --state-dirs="env/dev"
	@echo
	bin/tfmod dependents --debug --base-dir "testdata" --module-dirs="terraform/module/foo" --format=json

.PHONY: test
test: lint ## test
	go test -short ./...

.PHONY: test-all
test-all: lint ## test all
	go test -race -shuffle=on ./...

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

# Targets: Release
.PHONY: release
release: release/run ## Start release process
