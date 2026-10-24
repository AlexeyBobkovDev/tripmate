#!/usr/bin/env bash

die() {
	printf '%s\n' "$*" >&2
	exit 1
}
