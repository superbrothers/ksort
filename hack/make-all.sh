#!/usr/bin/env bash

set -e -o pipefail

[[ -n "$DEBUG" ]] && set -x

./hack/make-binaries.sh
./hack/make-archives.sh
# vim: ai ts=2 sw=2 et sts=2 ft=sh
