#!/usr/bin/env bash

SCRIPT_DIR="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck source=scripts/lib/die.sh
. "$SCRIPT_DIR/lib/die.sh"

for service in services/*; do
	[ -d "$service" ] || continue

	SERVICE_NAME=${service#services/}
	IMAGE="$SERVICE_NAME-test-docker-build:test-version"
	DOCKERFILE_PATH="cmd/$SERVICE_NAME/Dockerfile"
	(
		cd "$service" || die "can not cd to service=$service"
		[[ -f "$DOCKERFILE_PATH" ]] || die "Dockerfile not found $DOCKERFILE_PATH"
		docker build -t "$IMAGE" -f "$DOCKERFILE_PATH" . || exit 1
		docker rmi "$IMAGE" || exit 1
	) || exit 1
done
