#!/usr/bin/env bash

for script in scripts/*; do
	shellcheck -x "$script"
done
