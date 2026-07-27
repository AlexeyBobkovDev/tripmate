#!/usr/bin/env bash

set -euo pipefail

FILES=("$@")

test -z "$(gofmt -l "${FILES[@]}")" || echo "$(gofmt -l "${FILES[@]}")"
test -z "$(goimports -l "${FILES[@]}")" || echo "$(goimports -l "${FILES[@]}")"
test -z "$(gofumpt -l "${FILES[@]}")" || echo "$(gofumpt -l "${FILES[@]}")"
test -z "$(
	gci diff \
		--custom-order \
		-s standard \
		-s default \
		-s "prefix(github.com/AlexeyBobkovDev/tripmate)" \
		"${FILES[@]}"
)" || {
	echo "$(gci diff \
		--custom-order \
		-s standard \
		-s default \
		-s "prefix(github.com/AlexeyBobkovDev/tripmate)" \
		"${FILES[@]}")"
	exit 1
}

for service in services/*; do
	echo "$service"
	(
		cd "$service"

		golangci-lint run ./...
	)
done
