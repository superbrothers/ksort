#!/usr/bin/env bash

set -e -o pipefail

[[ -n "$DEBUG" ]] && set -x

SCRIPTDIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd "${SCRIPTDIR}/.."

bin_dir="out/bin"
if [[ ! -d "${bin_dir}" ]]; then
    echo "Binaries are not built (${bin_dir}), run make-binaries.sh" >&2
    exit 1
fi

for osarch in $(ls "$bin_dir"); do
  archive_file="out/ksort-${osarch}.zip"
  echo "Creating ${archive_file}" >&2
  ( \
    cd "${bin_dir}/${osarch}/" && \
    cp ../../../LICENSE.txt . && \
    cp ../../../README.md . && \
    zip -r "../../../${archive_file}" * \
  )

  tar_sumfile="${archive_file}.sha256"
  tar_checksum="$(shasum -a 256 "$archive_file" | awk '{print $1}')"
  echo "${archive_file}: checksum: ${tar_checksum}" >&2
  echo "$tar_checksum" >"${tar_sumfile}"
  echo "Written ${tar_sumfile}"
done

# Copy and process ksort manifest
cp hack/sort-by-kind.yaml out/sort-by-kind.yaml
git_tag="${TAG_NAME:-$(git describe --tags --dirty --always)}"
sed -i \
  -e "s/KSORT_DARWIN_ZIP_CHECKSUM/$(cat out/ksort-darwin-amd64.zip.sha256)/g" \
  -e "s/KSORT_LINUX_ZIP_CHECKSUM/$(cat out/ksort-linux-amd64.zip.sha256)/g" \
  -e "s/KSORT_VERSION/${git_tag}/g" \
  out/sort-by-kind.yaml
echo "Written out/sort-by-kind.yaml" >&2
# vim: ai ts=2 sw=2 et sts=2 ft=sh