.DEFAULT_GOAL := help

# Global vars
export SYS_GO=$(shell which go)
export SYS_GOFMT=$(shell which gofmt)
export SYS_DOCKER=$(shell which docker)

export BINARY_DIR=bin
export BINARY_NAME=owl_clerk_bot

.PHONY: run.bin
## Run & build bin
build:
	@echo 'Start build'
	@$(SYS_GO) build -o ./$(BINARY_DIR)/$(BINARY_NAME) cmd/app/main.go

.PHONY: build.bin
## Build bin file from go
build.bin: cmd/app/main.go
	$(SYS_GO) mod download && CGO_ENABLED=0 $(SYS_GO) build -o ./$(BINARY_DIR)/$(BINARY_NAME) cmd/app/main.go

.PHONY: run.cmd
## Run as go run cmd/app/main.go
run.cmd: cmd/app/main.go
	$(SYS_GO) run cmd/app/main.go

.PHONY: test
## Run test a make cover-file
test:
	@go test ./... -coverprofile cover.out

.PHONY: test.cover
## Run test coverage
test.cover: test
	$(SYS_GO) tool cover -html=cover.out

.PHONY: help
## Show this help message
help:
	@echo "$$(tput bold)Available rules:$$(tput sgr0)"
	@echo
	@sed -n -e "/^## / { \
		h; \
		s/.*//; \
		:doc" \
		-e "H; \
		n; \
		s/^## //; \
		t doc" \
		-e "s/:.*//; \
		G; \
		s/\\n## /---/; \
		s/\\n/ /g; \
		p; \
	}" ${MAKEFILE_LIST} \
	| LC_ALL='C' sort --ignore-case \
	| awk -F '---' \
		-v ncol=$$(tput cols) \
		-v indent=19 \
		-v col_on="$$(tput setaf 6)" \
		-v col_off="$$(tput sgr0)" \
	'{ \
		printf "%s%*s%s ", col_on, -indent, $$1, col_off; \
		n = split($$2, words, " "); \
		line_length = ncol - indent; \
		for (i = 1; i <= n; i++) { \
			line_length -= length(words[i]) + 1; \
			if (line_length <= 0) { \
				line_length = ncol - indent - length(words[i]) - 1; \
				printf "\n%*s ", -indent, " "; \
			} \
			printf "%s ", words[i]; \
		} \
		printf "\n"; \
	}' \
	| more $(shell test $(shell uname) == Darwin && echo '--no-init --raw-control-chars')
