#!/usr/bin/env bash

set -e -o pipefail

[[ -n "$DEBUG" ]] && set -x

krew_version="v0.2.1"

plugin_manifest_file="out/sort-manifests.yaml"
if [[ ! -f "$plugin_manifest_file" ]]; then
    echo "Plugin manifest file does not exist (${plugin_manifest_file}), run make-archives.sh" >&2
    exit 1
fi

krew_command="$HOME/.krew/bin/kubectl-krew"
if [[ ! -f "$krew_command" ]]; then
  (
    set -x; cd "$(mktemp -d)" &&
    curl -fsSLO "https://storage.googleapis.com/krew/${krew_version}/krew.{tar.gz,yaml}" &&
    tar zxvf krew.tar.gz &&
    ./krew-"$(uname | tr '[:upper:]' '[:lower:]')_amd64" install \
      --manifest=krew.yaml --archive=krew.tar.gz
  )
fi

export KREW_ROOT="$(mktemp -d)"
trap "rm -rf $KREW_ROOT" EXIT

"$HOME/.krew/bin/kubectl-krew" install \
    --manifest "$plugin_manifest_file" \
    --archive out/ksort-$(uname -s | tr "[:upper:]" "[:lower:]")-amd64.zip
# vim: ai ts=2 sw=2 et sts=2 ft=sh
