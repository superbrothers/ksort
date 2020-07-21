#!/usr/bin/env bash

set -e -o pipefail

[[ -n "$DEBUG" ]] && set -x

SCRIPTDIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd "${SCRIPTDIR}/.."

DEFAULT_OSARCH="darwin/amd64 linux/amd64 windows/amd64"
version_pkg="github.com/superbrothers/ksort"

rm -rf out/

# Builds
echo >&2 "Building binaries for: ${OSARCH:-$DEFAULT_OSARCH}"
git_rev="${SHORT_SHA:-$(git rev-parse --short HEAD)}"
git_tag="${TAG_NAME:=$(git describe --tags --abbrev=0 --exact-match 2>/dev/null ||:)}"
if [[ -z "$git_tag" ]]; then
  git_tag="$git_rev"
fi
if [[ -n "$(git status --porcelain)" ]]; then
  git_tree_state="dirty"
else
  git_tree_state="clean"
fi
build_date="$(date -u +'%Y-%m-%dT%H:%M:%SZ')"

env CGO_ENABLED=0 GO111MODULE=on gox -osarch="${OSARCH:-$DEFAULT_OSARCH}" \
  -tags netgo \
  -ldflags="-X ${version_pkg}.GitCommit=${git_rev} \
    -X ${version_pkg}.GitVersion=${git_tag} \
    -X ${version_pkg}.GitTreeState=${git_tree_state} \
    -X ${version_pkg}.BuildDate=${build_date}" \
  -output="out/bin/{{.OS}}-{{.Arch}}/ksort" \
  ./cmd/ksort
# vim: ai ts=2 sw=2 et sts=2 ft=sh
