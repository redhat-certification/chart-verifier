#!/bin/sh
set -eu

# Usage (after user approval): sh install-chart-verifier.sh [latest|TAG] [DESTINATION_DIR]
version="${1:-latest}"
dest="${2:-.}"

for dependency in gh curl tar mkdir mktemp chmod ln; do
  if ! command -v "${dependency}" >/dev/null 2>&1; then
    printf 'Missing dependency: %s\n' "${dependency}" >&2
    exit 1
  fi
done

binary="${dest}/chart-verifier"
if [ -e "${binary}" ] || [ -L "${binary}" ]; then
  echo "Refusing to overwrite ${binary}" >&2
  exit 1
fi

endpoint="repos/redhat-certification/chart-verifier/releases/tags/${version}"
if [ "${version}" = latest ]; then
  endpoint=repos/redhat-certification/chart-verifier/releases/latest
fi
url="$(gh api --hostname github.com "${endpoint}" --jq '
  [.assets[] | select(.name | startswith("chart-verifier") and endswith(".tgz"))]
  | if length == 1 then
      .[0].browser_download_url
    else error("Expected exactly one chart-verifier tarball") end')"
case "${url}" in
  https://github.com/redhat-certification/chart-verifier/releases/download/*) ;;
  *)
    echo "Unexpected asset URL" >&2
    exit 1
    ;;
esac

mkdir -p "${dest}"
tmp="$(mktemp -d "${dest}/.chart-verifier.XXXXXX")"
curl -fSL --proto '=https' -o "${tmp}/release.tgz" "${url}"
tar -xzOf "${tmp}/release.tgz" chart-verifier > "${tmp}/chart-verifier"
chmod +x "${tmp}/chart-verifier"
# A hard link installs atomically without overwriting an existing file.
ln -T "${tmp}/chart-verifier" "${binary}"
"${binary}" version
"${binary}" --help
printf 'Installed %s\n' "${binary}"
