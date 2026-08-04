#!/usr/bin/env bash

for service in services/*; do
	[ -d "$service" ] || continue

	SERVICE_NAME=${service#services/}
	IMAGE="$SERVICE_NAME-test-docker-build:test-version"
	DOCKERFILE_PATH="cmd/$SERVICE_NAME/Dockerfile"
	(
		cd "$service" || dye "can not cd to service=$service"
		docker build -t "$IMAGE" -f "$DOCKERFILE_PATH" .
		docker rmi "$IMAGE"
	)
done
