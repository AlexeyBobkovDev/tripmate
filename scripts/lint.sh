#!/usr/bin/env bash

SCRIPT_DIR="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck source=scripts/lib/die.sh
. "$SCRIPT_DIR/lib/die.sh"

for service in services/*; do
    [ -d "$service" ] || continue

    (
        cd "$service" || die "can not cd to service=$service"

        golangci-lint run ./...
    ) || exit 1
done
