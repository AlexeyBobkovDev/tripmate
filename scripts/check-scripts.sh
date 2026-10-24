#!/usr/bin/env bash

SCRIPT_DIR="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck source=scripts/lib/die.sh
. "$SCRIPT_DIR/lib/die.sh"

SH_FILES=()

for file in "$@"; do
	[[ "$file" == *.sh ]] && SH_FILES+=("$file")
done

if [[ ${#SH_FILES[@]} != 0 ]]; then
	shellcheck -x "${SH_FILES[@]}" || exit 1
else
	mapfile -d '' -t ALL_SH_FILES < <(find scripts -type f -name '*.sh' -print0)
	shellcheck -x "${ALL_SH_FILES[@]}" || exit 1
fi
