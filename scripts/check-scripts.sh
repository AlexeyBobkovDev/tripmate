#!/usr/bin/env bash

output="$(shellcheck -x scripts/*.sh)"

test -z "$output" || echo "$output"
