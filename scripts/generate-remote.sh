#!/usr/bin/env bash
#
# Reference script for building a TGENV_REMOTE-compatible terragrunt mirror.
#
# Produces the file layout tgswitch expects when TGENV_REMOTE is set:
#   <output_dir>/versions.json
#   <output_dir>/vX.Y.Z/terragrunt_<os>_<arch>
#
# Upload the contents of <output_dir> to whatever you point TGENV_REMOTE at
# (S3, GCS, an internal artifact repo, a plain static file server, ...).
#
# Usage:
#   VERSIONS="1.1.0 1.1.1" ./generate-remote.sh [output_dir]
#
# Env vars:
#   VERSIONS  Required. Space-separated terragrunt versions to mirror (no "v" prefix).
#   ENDPOINT  Source to download release binaries from. Default: the public
#             gruntwork-io/terragrunt GitHub releases page.
#   GOOS      Target OS for the binary name. Default: linux.
#   GOARCH    Target arch for the binary name. Default: amd64.
#
# Requires: curl, jq

set -euo pipefail

ENDPOINT="${ENDPOINT:-https://github.com/gruntwork-io/terragrunt}"
GOOS="${GOOS:-linux}"
GOARCH="${GOARCH:-amd64}"
OUT_DIR="${1:-terragrunt-releases}"

if [[ -z "${VERSIONS:-}" ]]; then
  echo "Set VERSIONS to a space-separated list of terragrunt versions to mirror, e.g.:" >&2
  echo "  VERSIONS=\"1.1.0 1.1.1\" $0" >&2
  exit 1
fi

download_artifact() {
  local version="$1" dest="$2"
  local binary="terragrunt_${GOOS}_${GOARCH}"
  local url="${ENDPOINT}/releases/download/v${version}/${binary}"

  echo "Downloading terragrunt ${version} (${GOOS}/${GOARCH})..."
  curl -LsS --fail -o "${dest}/${binary}" "$url"
}

mkdir -p "$OUT_DIR"

versions_json='{"Versions": []}'

for version in $VERSIONS; do
  version_dir="${OUT_DIR}/v${version}"
  mkdir -p "$version_dir"
  download_artifact "$version" "$version_dir"
  versions_json=$(jq --arg v "$version" '.Versions += [$v]' <<<"$versions_json")
done

jq . <<<"$versions_json" > "${OUT_DIR}/versions.json"

echo
echo "Remote layout ready in ${OUT_DIR}/"
echo "Upload its contents to your TGENV_REMOTE endpoint."
