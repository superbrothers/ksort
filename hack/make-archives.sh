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

# Create archive files in zip and tar.gz format for each arch
for osarch in $(ls "$bin_dir"); do
  for format in zip tar.gz; do
    archive_file="out/ksort-${osarch}.${format}"
    echo "Creating ${archive_file}" >&2
    ( \
      cd "${bin_dir}/${osarch}/" && \
      cp ../../../LICENSE.txt . && \
      cp ../../../README.md . && \
      [[ "$format" == "zip" ]] && \
        zip -r "../../../${archive_file}" * \
      || \
        tar zcvf "../../../${archive_file}" * \
    )

    archive_sumfile="${archive_file}.sha256"
    archive_checksum="$(shasum -a 256 "$archive_file" | awk '{print $1}')"
    echo "${archive_file}: checksum: ${archive_checksum}" >&2
    echo "$archive_checksum" >"${archive_sumfile}"
    echo "Written ${archive_sumfile}"
  done
done

# Copy and process ksort manifest
cp hack/sort-manifests.yaml out/sort-manifests.yaml
git_tag="${TRAVIS_TAG:-$(git describe --tags --dirty --always)}"
sed -i \
  -e "s/KSORT_DARWIN_ZIP_CHECKSUM/$(cat out/ksort-darwin-amd64.zip.sha256)/g" \
  -e "s/KSORT_LINUX_ZIP_CHECKSUM/$(cat out/ksort-linux-amd64.zip.sha256)/g" \
  -e "s/KSORT_WINDOWS_ZIP_CHECKSUM/$(cat out/ksort-windows-amd64.zip.sha256)/g" \
  -e "s/KSORT_VERSION/${git_tag}/g" \
  out/sort-manifests.yaml
echo "Written out/sort-manifests.yaml" >&2
# vim: ai ts=2 sw=2 et sts=2 ft=sh
