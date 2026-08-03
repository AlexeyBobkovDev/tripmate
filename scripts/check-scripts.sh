#!/usr/bin/env bash

find scripts -name "*.sh" -type f -exec shellcheck -x {} +
