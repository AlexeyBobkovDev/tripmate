set -euo pipefail

for service in services/*; do
    [ -d "$service" ] || continue

    (
        cd "$service" || exit 1

        golangci-lint run ./...
    )
done
