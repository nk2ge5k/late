REPO = late
ROOT_DIR = $(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))

GO ?= $(shell which go)
PYTHON ?= $(shell which python3)
GOLANGCI_LINT_FORMAT ?= "colored-line-number"

include Makefile.tools

build: gen ## Build for production
	@mkdir -p $(ROOT_DIR)/build/production
	$(GO) build -o $(ROOT_DIR)/build/production/late $(ROOT_DIR)/cmd/late

build-test: gen ## Build for local testing
	@mkdir -p $(ROOT_DIR)/build/development
	$(GO) build -o $(ROOT_DIR)/build/development/late $(ROOT_DIR)/cmd/late


##################################### TEST #####################################

test: test-go test-python ## Run all tests
	@echo ""

test-go: gen ## Run all unit tests
	@$(GO) test -v -race ./...

test-python: $(VENV) gen build-test ## Run all python integration tests
	@PYTHON=$(VENV_PYTHON) $(ROOT_DIR)/tests/run-testsuite.sh $(ROOT_DIR)

################################### GENERATE ###################################

gen: gen-go gen-proto ## Run code generation
	@echo ""

gen-go: ## Generate files with go-generate
	@$(GO) generate ./...

gen-proto: $(VENV) $(BUF) $(PROTOC_GEN_GO) $(PROTOC_GEN_GO_GRPC) $(PROTOC_GEN_GO_GRPC_GATEWAY) ## Generate code with protoc plugins
	@$(TMP_BIN)/buf generate \
		--template "$(shell \
				ROOT_DIR=$(ROOT_DIR) \
				REPO=$(REPO) \
				envsubst < $(ROOT_DIR)/buf.gen.tmpl.json)" || true

#################################### FORMAT ####################################

format: format-proto format-yaml format-python format-go ## Format service files
	@echo ""

format-go: $(GOIMPORTS) ## Format service Go files
	@$(TMP_BIN)/goimports -local=$(REPO) -w \
		$(shell ls -d $(ROOT_DIR)/*/ | grep -v -e vendor -e .tmp -e .venv) \
		$(shell ls $(ROOT_DIR) | grep .go)

format-proto: $(BUF) ## Format Protobuf files
	@$(TMP_BIN)/buf format -w --config $(ROOT_DIR)/buf.yaml

format-yaml: $(YAMLFMT) ## Format YAML files
	@$(TMP_BIN)/yamlfmt -conf $(ROOT_DIR)/.yamlfmt .

format-python: $(VENV) ## Format Python files
	@$(PYTHON_VENV_DIR)/bin/isort . 2>/dev/null > /dev/null
	@$(PYTHON_VENV_DIR)/bin/black . 2>/dev/null > /dev/null

##################################### LINT #####################################

lint: lint-proto lint-go ## Run every available linter
	@echo ""

lint-go: $(GOLANGCI_LINT) gen ## Run linting on Go files
	@$(TMP_BIN)/golangci-lint run \
		-c $(ROOT_DIR)/.golangci.yaml \
		--out-format=$(GOLANGCI_LINT_FORMAT)

lint-proto: $(BUF) ## Run linting on Protobuf files
	@$(TMP_BIN)/buf lint -v --config $(ROOT_DIR)/buf.yaml

##################################### HELP #####################################

help: ## Show this help
	@echo "\nSpecify a command. The choices are:\n"
	@grep -hE '^[0-9a-zA-Z_-]+:.*?## .*$$' ${MAKEFILE_LIST} \
		| awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[0;36m%-20s\033[m %s\n", $$1, $$2}'
	@echo ""
.PHONY: help

clean: ## Clean temporary files and directories
	git clean -xdf
	find . -type d -name __pycache__ -print | xargs rm -rf
	find . -type d -name '.pytest_cache' -print | xargs rm -rf

.DEFAULT_GOAL := help
