.DEFAULT_GOAL := test

# Clean up
clean:
	@rm -fR ./vendor/ ./cover.*
.PHONY: clean

# Download project dependencies
configure:
	@GO111MODULE=on go mod download
.PHONY: configure

# Run tests and generates html coverage file
cover: test
	@go tool cover -html=./cover.out -o ./cover.html
	@test -f ./cover.out && rm ./cover.out;
.PHONY: cover

# Format all go files
fmt:
	@gofmt -s -w -l $(shell go list -f {{.Dir}} ./...)
.PHONY: fmt

# Run tests
test:
	@go test -v -coverprofile=./cover.out -covermode=atomic $(shell go list ./...)
.PHONY: test
