#!/usr/bin/env bash

SCRIPT_DIR="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck source=scripts/lib/die.sh
. "$SCRIPT_DIR/lib/die.sh"

FILES=("$@")
if (( ${#FILES[@]} == 0 )); then
	echo "One or more files must be provided"
	exit 1
fi

unformatted="$(gofmt -l "${FILES[@]}")"  || exit 1
if [[ -n "$unformatted" ]]; then
	printf '%s\n' "$unformatted"
	exit 1
fi
unformatted="$(goimports -l "${FILES[@]}")"  || exit 1
if [[ -n "$unformatted" ]]; then
	printf '%s\n' "$unformatted"
	exit 1
fi
unformatted="$(gofumpt -l "${FILES[@]}")"  || exit 1
if [[ -n "$unformatted" ]]; then
	printf '%s\n' "$unformatted"
	exit 1
fi

unformatted="$(gci diff \
		--custom-order \
		-s standard \
		-s default \
		-s "prefix(github.com/AlexeyBobkovDev/tripmate)" \
		"${FILES[@]}" </dev/null
)" || exit 1
if [[ -n "$unformatted" ]]; then
	printf '%s\n' "$unformatted"
	exit 1
fi

for service in services/*; do
	[ -d "$service" ] || continue
	(
		cd "$service" || die "can not cd to service=$service"

		golangci-lint run ./...
	) || exit 1
done
