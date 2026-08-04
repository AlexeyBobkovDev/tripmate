#!/usr/bin/env bash

FILES=("$@")

test -z "$(gofmt -l "${FILES[@]}")" || gofmt -l "${FILES[@]}"
test -z "$(goimports -l "${FILES[@]}")" || goimports -l "${FILES[@]}"
test -z "$(gofumpt -l "${FILES[@]}")" || gofumpt -l "${FILES[@]}"
test -z "$(
	gci diff \
		--custom-order \
		-s standard \
		-s default \
		-s "prefix(github.com/AlexeyBobkovDev/tripmate)" \
		"${FILES[@]}"
)" || {
	gci diff \
		--custom-order \
		-s standard \
		-s default \
		-s "prefix(github.com/AlexeyBobkovDev/tripmate)" \
		"${FILES[@]}"
	exit 1
}

for service in services/*; do
	echo "$service"
	(
		cd "$service" || die "can not cd to service=$service"

		golangci-lint run ./...
	)
done
