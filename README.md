# Weave Blog Tutorial

This repo demonstrates implementation of a blog application using [Weave](https://github.com/iov-one/weave).

## Documentation

Documentation of this repos is at [Weave Tutorial](https://docs.iov.one/docs/weave-tutorial/overview). You can follow the documentation while reading the code.

## Installing dependencies

### Requirements

- [golang 1.11.4+](https://golang.org/doc/install)
- [tendermint 0.31.5](https://github.com/tendermint/tendermint/tree/v0.31.5)
  - [Installation](https://github.com/tendermint/tendermint/blob/master/docs/introduction/install.md)
- [weave](https://github.com/iov-one/weave)
  - `go get github.com/iov-one/weave`
- [docker](https://docs.docker.com/install/)

**Important**: At IOV we use [go modules](https://github.com/golang/go/wiki/Modules). You can append `export GO111MODULE=on` to `.rc` file of your favorite shell(`zsh, bash, fish, etc.`)
