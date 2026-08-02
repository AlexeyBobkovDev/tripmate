#!/usr/bin/env bash

for script in scripts/*; do
	shellcheck "$script"
done
