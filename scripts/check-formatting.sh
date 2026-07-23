#!/usr/bin/env bash

FILES=("$@")

test -z "$(gofmt -l ${FILES})"
test -z "$(goimports -l ${FILES})"
test -z "$(gofumpt -l ${FILES})"
test -z "$(gci diff ${FILES})"
test -z "$(golangci-lint run)"
