# ======================================================
# Global variables
# ======================================================

GO_VERSION							= 1.24
GOLANGCI_LINT_VERSION		= v2.1.5
GCI_VERSION							= v0.13.6
GOFUMPT_VERSION					= v0.8.0

BIN_DIR					= $(CURDIR)/bin
GOLANGCI_LINT		= $(BIN_DIR)/golangci-lint
GCI							= $(BIN_DIR)/gci
GOFUMPT					= $(BIN_DIR)/gofumpt
GOOSE						= $(BIN_DIR)/goose

DOCKER_DIR		= $(CURDIR)/build/dev
DOCKER_YML		= $(DOCKER_DIR)/docker-compose.yml
ENV_NAME			= "wbtech-l0"

MODULES				= cmd internal
MIGRATE_DIR		= $(CURDIR)/migrations

-include ./build/dev/.env


# ======================================================
# Docker composer
# ======================================================

.PHONY: compose-up
compose-up:
	docker-compose -p ${ENV_NAME} -f ${DOCKER_YML} up -d

.PHONY: compose-down
compose-down: ## terminate local env
	docker-compose -p ${ENV_NAME} -f ${DOCKER_YML} stop

.PHONY: compose-rm
compose-rm: ## remove local env
	docker-compose -p ${ENV_NAME} -f ${DOCKER_YML} rm -fvs

.PHONY: compose-rs
compose-rs: ## remove previously and start new local env
	make compose-rm
	make compose-up

.PHONY: compose-app
compose-app:
	docker-compose -p ${ENV_NAME} -f ${DOCKER_YML} build order-app


# ======================================================
# Migration
# ======================================================

.PHONY: install-goose
install-goose:
	@echo "Installing goose if missing..."
	@[ -f $(GOOSE) ] || GOBIN=$(BIN_DIR) go install github.com/pressly/goose/v3/cmd/goose@latest


# ======================================================
# Formatter installation and formatting
# ======================================================

.PHONY: install-formatters
install-formatters:
	@echo "Installing formatters if missing..."
	@[ -f $(GOFUMPT) ] || GOBIN=$(BIN_DIR) go install mvdan.cc/gofumpt@$(GOFUMPT_VERSION)
	@[ -f $(GCI) ] || GOBIN=$(BIN_DIR) go install github.com/daixiang0/gci@$(GCI_VERSION)

.PHONY: format
format: install-formatters
	@echo "Formatting code with gofumpt..."
	@for module in $(MODULES); do \
		if [ -d $$module ]; then \
			find $$module -type f -name '*.go' ! -path '*/mocks/*' -exec $(GOFUMPT) -extra -w {} +; \
		fi \
	done
	@echo "Sorting imports with gci..."
	@for module in $(MODULES); do \
		if [ -d $$module ]; then \
			find $$module -type f -name '*.go' ! -path '*/mocks/*' -exec $(GCI) write -s standard \
			-s default -s "prefix(github.com/AndrejDubinin/wbtech-l0)" {} +; \
		fi \
	done


# ======================================================
# Linting
# ======================================================

.PHONY: install-golangci-lint
install-golangci-lint:
	@echo "Installing golangci-lint if missing..."
	@[ -f $(GOLANGCI_LINT) ] || { mkdir -p $(BIN_DIR); GOBIN=$(BIN_DIR) go install \
	github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION); }

.PHONY: lint
lint: install-golangci-lint
	@echo "Running golangci-lint on all modules..."
	@set -e; ERR=0; \
	for mod in $(MODULES); do \
		if [ -d $$mod ]; then \
			$(GOLANGCI_LINT) run $$mod/... --config=.golangci.yml || ERR=1; \
		fi \
	done; exit $$ERR


# ======================================================
# Local execution
# ======================================================

.PHONY: run-all
run-all:
	@DB_CONN=postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@localhost:${POSTGRES_HOST_PORT}/${POSTGRES_DB} \
	CACHE_CAPACITY=100 go run ./cmd/app | jq -c 'del(.ts) | del(.stacktrace)'
