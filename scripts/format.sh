#!/usr/bin/env bash

SCRIPT_DIR="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck source=scripts/lib/die.sh
. "$SCRIPT_DIR/lib/die.sh"

GO_FILES=()
for file in "$@"; do
	[[ "$file" == *.go ]] && GO_FILES+=("$file")
done

if (( ${#GO_FILES[@]} > 0 )); then
    FILES=("${GO_FILES[@]}")
else
    mapfile -d '' -t FILES < <(find services -type f -name '*.go')
fi

for file in "${FILES[@]}"; do
    service=$(echo "$file" | cut -d/ -f1-2)
    [ -d "$service" ] || continue

    (
        cd "$service" || die "can not cd to service=$service"

        rel="${file#"$service"/}"

        goimports -w "$rel" || exit 1
        gofumpt -w "$rel" || exit 1
        gci write \
			--custom-order \
			-s standard \
			-s default \
			-s "prefix(github.com/AlexeyBobkovDev/tripmate)" \
			"$rel" || exit 1
    ) || exit 1
done
