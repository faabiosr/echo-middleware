.DEFAULT_GOAL := help

.PHONY: clean
clean: ## clean up files generate by coverage or go mod
	@rm -fR ./vendor/ ./cover.*

.PHONY: cover
cover: test ## run tests and generates the html coverage file
	@go tool cover -html=./cover.out -o ./cover.html
	@test -f ./cover.out && rm ./cover.out;

.PHONY: help
help: ## display help screen
	@echo "Usage: make <target>"
	@echo ""
	@sed \
		-e '/^[a-zA-Z0-9_\-]*:.*##/!d' \
		-e 's/:.*##\s*/:/' \
		-e 's/^\(.\+\):\(.*\)/$(shell tput setaf 6)\1$(shell tput sgr0):\2/' \
		$(MAKEFILE_LIST) | column -c2 -t -s :
	@echo ''

.PHONY: lint
lint: # golang linters (golangci-lint)
	@golangci-lint run ./...

.PHONY: test
test: ## run tests
	@go test -v -coverprofile=./cover.out -covermode=atomic $(shell go list ./...)
