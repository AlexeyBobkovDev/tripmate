#!/usr/bin/env bash

for service in services/*; do
    [ -d "$service" ] || continue

    (
        cd "$service" || dye "can not cd to service=$service"

        golangci-lint run ./...
    )
done
