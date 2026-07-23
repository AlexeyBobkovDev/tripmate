#!/usr/bin/env bash

set -euo pipefail

FILES=("$@")


test -z "$(gofmt -l ${FILES[@]})"
test -z "$(goimports -l ${FILES[@]})"
test -z "$(gofumpt -l ${FILES[@]})"
test -z "$(gci diff ${FILES[@]})"

for service in services/*; do
	(
		cd "$service"
		for file in ${FILES[@]}; do
			if [[ "$file" != "$service/"* ]]; then
				continue
			fi

			file="${file#$service/}"
			
			test -z "$(golangci-lint run "$file")"
		done
	)
done

