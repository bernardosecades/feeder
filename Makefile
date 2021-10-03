TEST_DB_VARS = ENV=test DB_PORT=5416 DB_USER=feeder DB_PASS=feeder DB_NAME=feeder DB_HOST=127.0.0.1
.PHONY: mod
## Install project dependencies using go mod. Usage 'make mod'
mod:
	@go mod tidy
	@go mod vendor
.PHONY: install
## Install dependency manually. Usage 'make install pkg=path[version/hash_commit]'
install:
	@go get $(pkg)
.PHONY: lint
## Run linter. Usage: 'make lint'
lint: ; $(info Linting feeder...)
	golangci-lint run --fix
.PHONY: run
## Run feeder service. Usage: 'make run'
run: ; $(info Starting feed sever...)
	go run ./cmd/feedersrv/.

.PHONY: test
## Run tests. Usage: 'make test' Options: path=./some-path/... [and/or] func=TestFunctionName
test: ; $(info running tests...) @
	@if [ -z $(path) ]; then \
		path='./...'; \
	else \
		path=$(path); \
	fi; \
	if [ ! -d "coverage" ]; then \
		mkdir coverage; \
	fi; \
	if [ -z $(func) ]; then \
		$(TEST_DB_VARS) go test -v -failfast -covermode=count -coverprofile=./coverage/coverage.out $$path; \
	else \
		$(TEST_DB_VARS) go test -v -failfast -covermode=count -coverprofile=./coverage/coverage.out -run $$func $$path; \
	fi; \

# COLORS
GREEN  := $(shell tput -Txterm setaf 2)
YELLOW := $(shell tput -Txterm setaf 3)
WHITE  := $(shell tput -Txterm setaf 7)
RESET  := $(shell tput -Txterm sgr0)
TARGET_MAX_CHAR_NUM=20
## Show help
help:
	@echo ''
	@echo 'Usage:'
	@echo '  ${YELLOW}make${RESET} ${GREEN}<target>${RESET}'
	@echo ''
	@echo 'Targets:'
	@awk '/^[a-zA-Z\-\_0-9]+:/ { \
		helpMessage = match(lastLine, /^## (.*)/); \
		if (helpMessage) { \
			helpCommand = substr($$1, 0, index($$1, ":")); \
			helpMessage = substr(lastLine, RSTART + 3, RLENGTH); \
			printf "  ${YELLOW}%-$(TARGET_MAX_CHAR_NUM)s${RESET} ${GREEN}%s${RESET}\n", helpCommand, helpMessage; \
		} \
	} \
	{ lastLine = $$0 }' $(MAKEFILE_LIST)
