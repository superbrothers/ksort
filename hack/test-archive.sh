#!/usr/bin/env bash

set -e -o pipefail

[[ -n "$DEBUG" ]] && set -x


plugin_manifest_file="out/sort-by-kind.yaml"
if [[ ! -f "$plugin_manifest_file" ]]; then
    echo "Plugin manifest file does not exist (${plugin_manifest_file}), run make-archives.sh" >&2
    exit 1
fi

export KREW_ROOT="$(mktemp -d)"
trap "rm -rf $KREW_ROOT" EXIT


"$HOME/.krew/bin/kubectl-krew" install \
    --manifest out/sort-by-kind.yaml \
    --archive out/ksort-$(uname -s | tr "[:upper:]" "[:lower:]")-amd64.zip
# vim: ai ts=2 sw=2 et sts=2 ft=sh
