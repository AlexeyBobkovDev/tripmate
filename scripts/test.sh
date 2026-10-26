#!/usr/bin/env bash

SCRIPT_DIR="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck source=scripts/lib/die.sh
. "$SCRIPT_DIR/lib/die.sh"

GO_FILES=()

if (( "$#" == 0 )); then
	mapfile -d '' -t GO_FILES < <(find services -type f -name '*.go' -print0)
else
	for file in "$@"; do
		[[ "$file" == *.go ]] && GO_FILES+=("$file")
	done
fi

declare -A SERVICES
for file in "${GO_FILES[@]}"; do
	service=$(echo "$file" | cut -d/ -f1-2)
	[ -d "$service" ] || continue
	SERVICES["$service"]=true
done

for service in "${!SERVICES[@]}"; do
	[ -d "$service" ] || continue

	(
		cd "$service" || die "can not cd to service=$service"
		go test ./... || exit 1

		if [ "${TEST_INTEGRATION:-0}" == 1 ]; then
			LOGGER_FOLDER=../../../out/logs go test -tags=integration -race -count=1 -v ./... || exit 1
		fi
		if [ "${TEST_E2E:-0}" == 1 ]; then
			LOGGER_FOLDER=../../../out/logs go test -tags=e2e -race -count=1 -v ./... || exit 1
		fi

	) || exit 1
done
