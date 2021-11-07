#!/usr/bin/env bash

# ARCH_TYPE_1
# examples:
#   - x86_64
#   - arm64

ARCH=$(uname -m)
export ARCH_TYPE_1="${ARCH}"

# ARCH_TYPE_2
# examples:
#   - amd64
#   - arm64

ARCH=$(uname -m)
if [ "${ARCH}" = "x86_64" ]; then
    export ARCH_TYPE_2="amd64"
else
    export ARCH_TYPE_2=${ARCH}
fi

# ARCH_TYPE_3
ARCH=$(uname -m)
if [ "${ARCH}" = "arm64" ]; then
    export ARCH_TYPE_3="arm64v8"
elif [ "${ARCH}" = "x86_64" ]; then
    export ARCH_TYPE_3="amd64"
else
    export ARCH_TYPE_3=${ARCH}
fi
