# Contributing

## Setup your machine

`echo-middleware` is written in [Go](https://golang.org/).

Prerequisites:

* `make`
* [Go 1.21+](https://golang.org/doc/install)
* [GolangCI-Lint](https://github.com/golangci/golangci-lint#install)

Clone `echo-middleware` from source:

```sh
$ git clone git@github.com:faabiosr/echo-middleware.git
$ cd echo-middleware
```

A good way of making sure everything is all right is running the test suite:
```console
$ make test
```

## Create a commit

### Commit message format and code base

Commit messages should be well formatted (please use [gofumpt](https://github.com/mvdan/gofumpt)), and must follow the [Conventional Commits](https://www.conventionalcommits.org) specification, this pattern was defined to help with new version releases.

Start your commit message with the type. Choose one of the following:
`feat`, `fix`, `docs`, `style`, `refactor`, `perf`, `test`, `chore`, `revert`, `add`, `remove`, `move`, `bump`, `update`, `release`

After a colon, you should give the message a title, starting with uppercase and ending without a dot.
Keep the width of the text at 72 chars.
The title must be followed with a newline, then a more detailed description.

Please reference any GitHub issues on the last line of the commit message (e.g. `See #123`, `Closes #123`, `Fixes #123`).

An example:

```
docs: add example for --release-notes flag

I added an example to the docs of the `--release-notes` flag to make
the usage more clear.  The example is an realistic use case and might
help others to generate their own changelog.

See #284
```

## Submit a pull request

Push your branch to your `echo-middleware` fork and open a pull request against the main branch.
