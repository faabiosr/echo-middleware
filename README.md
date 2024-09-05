# Echo Framework - Middlewares

[![Build Status](https://img.shields.io/github/actions/workflow/status/faabiosr/echo-middleware/test.yml?logo=github&style=flat-square)](https://github.com/faabiosr/echo-middleware/actions?query=workflow:test)
[![Codecov branch](https://img.shields.io/codecov/c/github/faabiosr/echo-middleware/master.svg?style=flat-square)](https://codecov.io/gh/faabiosr/echo-middleware)
[![Go Reference](https://pkg.go.dev/badge/github.com/faabiosr/echo-middleware.svg)](https://pkg.go.dev/github.com/faabiosr/echo-middleware)
[![Go Report Card](https://goreportcard.com/badge/github.com/faabiosr/echo-middleware?style=flat-square)](https://goreportcard.com/report/github.com/faabiosr/echo-middleware)
[![License](https://img.shields.io/badge/License-MIT-blue.svg?style=flat-square)](https://github.com/faabiosr/echo-middleware/blob/master/LICENSE)

## :tada: Overview
echo-middleware is a Go package that provides multiple middlewares for Echo Framework.

## :relaxed: Motivation
After writing these middlewares several times, it was decided to create a package with useful middlewares for echo.

## :dart: Installation

### Requirements
echo-middleware requires Go 1.21 or later and Echo Framework v4.

### How to use
Use go get.
```sh
$ go get github.com/faabiosr/echo-middleware
```

Then import the package into your own code:
```
import "github.com/faabiosr/echo-middleware"
```

## :toolbox: Development

### Requirements
- Install [Go](https://golang.org)
- Install [GolangCI-Lint](https://github.com/golangci/golangci-lint#install)

### Makefile
Please run `make help` to see all the available targets.

## :page_with_curl: License
This project is released under the MIT licence. See [LICENSE](https://github.com/faabiosr/echo-middleware/blob/master/LICENSE) for more details.
