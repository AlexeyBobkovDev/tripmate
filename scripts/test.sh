set -euo pipefail

GO_FILES=$(printf "%s\n" "$@" | grep -E '\.go$' || true)

if [ -n "$GO_FILES" ]; then
	FILES="$GO_FILES"
else
	FILES=$(find services -type f -name '*.go')
fi

for file in $FILES; do
	service=$(echo "$file" | cut -d/ -f1-2)
	[ -d "$service" ] || continue

	(
		cd "$service" || exit 1

		rel="${file#"$service"/}"
		go test ./...
		if [ ${TEST_INTEGRATION:-0} -eq 1 ]; then
			go test -tags=integration ./...
		fi
	)
done
