SHELL := /bin/bash
.SILENT:
.DEFAULT_GOAL := help

# Global vars
export SYS_GO=$(shell which go)
export SYS_GOFMT=$(shell which gofmt)
export BINARY_DIR=dist
export BINARY_NAME=owl_clerk_bot


.PHONY: run.cmd
## Run as go run cmd/app/main.go
run.cmd: cmd/app/main.go
	$(SYS_GO) run cmd/app/main.go

.PHONY: run.bin
## Run as binary
run.bin: build
	source .env && ./$(BINARY_DIR)/$(BINARY_NAME)

.PHONY: install-deps
## Install all requirements
install-deps: go.mod
	$(SYS_GO) mod tidy

.PHONY: build
## Build bot
build: cmd/app/main.go
	$(SYS_GO) mod download && CGO_ENABLED=0 $(SYS_GO) build -o ./$(BINARY_DIR)/$(BINARY_NAME) cmd/app/main.go

.PHONY: fmt
## Run go fmt
fmt:
	$(SYS_GOFMT) -s -w .

.PHONY: vet
## Run go vet ./...
vet:
	$(SYS_GO) vet ./...

.PHONY: clean
## Clean all artifacts
clean:
	rm -rf $(BINARY_DIR)

.PHONY: test
## Run all test
test:
	go test --short -coverprofile=cover.out -v ./...
	make test.coverage

#.PHONY: test.integration
### Run test integration
#test.integration:
#	docker run --rm -d -p 27019:27017 --name $$TEST_CONTAINER_NAME -e MONGODB_DATABASE=$$TEST_DB_NAME mongo:4.4-bionic
#
#	GIN_MODE=release go test -v ./tests/
#	docker stop $$TEST_CONTAINER_NAME

#.PHONY: test.coverage
### Run test coverage
#test.coverage:
#	go tool cover -func=cover.out

#.PHONY: swag
### Run swag
#swag:
#	swag init -g internal/app/app.go

.PHONY: lint
## Run golangci-lint
lint:
	golangci-lint -v run --out-format=colored-line-number

#.PHONY: gen
### Run mockgen
#gen:
#	mockgen -source=internal/service/service.go -destination=internal/service/mocks/mock.go
#	mockgen -source=internal/repository/repository.go -destination=internal/repository/mocks/mock.go

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
