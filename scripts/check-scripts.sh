#!/usr/bin/env bash

SH_FILES=$(printf "%s\n" "$@" | grep -E '\.sh$' || true)

if [[ "${#SH_FILES[@]}" != 0 ]]; then
	output="$(shellcheck -x "${SH_FILES[@]}")"
else
	output="$(shellcheck -x scripts/*.sh)"
fi

test -z "$output" || echo "$output"
