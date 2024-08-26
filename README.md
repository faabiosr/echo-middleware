# Echo Framework - Middlewares

[![Build Status](https://img.shields.io/github/actions/workflow/status/faabiosr/echo-middleware/test.yml?logo=github&style=flat-square)](https://github.com/faabiosr/echo-middleware/actions?query=workflow:test)
[![Codecov branch](https://img.shields.io/codecov/c/github/faabiosr/echo-middleware/master.svg?style=flat-square)](https://codecov.io/gh/faabiosr/echo-middleware)
[![Go Reference](https://pkg.go.dev/badge/github.com/faabiosr/echo-middleware.svg)](https://pkg.go.dev/github.com/faabiosr/echo-middleware)
[![Go Report Card](https://goreportcard.com/badge/github.com/faabiosr/echo-middleware?style=flat-square)](https://goreportcard.com/report/github.com/faabiosr/echo-middleware)
[![License](https://img.shields.io/badge/License-MIT-blue.svg?style=flat-square)](https://github.com/faabiosr/echo-middleware/blob/master/LICENSE)

## Description
echo-middleware is a Go package that provides multiple middlewares for Echo Framework.

## Requirements
echo-middleware requires Go 1.19 or later and Echo Framework v4.

## Installation
Use go get.
```sh
$ go get github.com/faabiosr/echo-middleware
```

Then import the package into your own code:
```
import "github.com/faabiosr/echo-middleware"
```

## Development

### Requirements
- Install [Go](https://golang.org)

### Makefile
```sh
# Clean up
$ make clean

# Download project dependencies
$ make configure

# Run tests and generates html coverage file
$ make cover

# Format all go files
$ make fmt

# Run tests
$make test
```

## License
This project is released under the MIT licence. See [LICENSE](https://github.com/faabiosr/echo-middleware/blob/master/LICENSE) for more details.
