set -euo pipefail

for service in services/*; do
	[ -d "$service" ] || continue

	SERVICE_NAME=${service#services/}
	IMAGE="$SERVICE_NAME-test-docker-build:test-version"
	DOCKERFILE_PATH="cmd/$SERVICE_NAME/Dockerfile"
	(
		cd $service
		docker build -t "$IMAGE" -f "$DOCKERFILE_PATH" .
		docker rmi "$IMAGE"
	)
done
