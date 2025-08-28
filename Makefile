BIN_DIR=${CURDIR}/bin
DOCKER_DIR=${CURDIR}/build/dev
DOCKER_YML=${DOCKER_DIR}/docker-compose.yml
ENV_NAME="wbtech-l0"
MIGRATE_DIR=${CURDIR}/migrations

-include ./build/dev/.env

.PHONY: compose-up
compose-up:
	docker-compose -p ${ENV_NAME} -f ${DOCKER_YML} up #-d

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

.PHONY: install-goose
install-goose:
	GOBIN=${BIN_DIR} go install github.com/pressly/goose/v3/cmd/goose@latest

.PHONY: run-all
run-all:
	@DB_CONN=postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@localhost:${POSTGRES_HOST_PORT}/${POSTGRES_DB} go run ./cmd/app
