#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(cd "${SCRIPT_DIR}/.." && pwd)"
BUILD_DIR="${ROOT_DIR}/build"
APP_NAME="ServerCommander"
MAIN_PKG="${ROOT_DIR}/src"

print_step() {
  echo
  echo "${1}"
}

require_tool() {
  if ! command -v "$1" >/dev/null 2>&1; then
    echo "Error: required tool '$1' is not available. Install it from $2." >&2
    exit 1
  fi
}

print_step "Checking prerequisites..."
require_tool go "https://go.dev/dl/"
require_tool git "https://git-scm.com/"

mkdir -p "${BUILD_DIR}"

build_target() {
  local os="$1"
  local arch="$2"
  local output="$3"
  local tags=( )

  if [[ -n "${GO_BUILD_TAGS:-}" ]]; then
    tags=(-tags "${GO_BUILD_TAGS}")
  fi

  print_step "Building ${os}/${arch}..."
  GOOS="${os}" GOARCH="${arch}" CGO_ENABLED=0 \
    go build "${tags[@]}" -o "${output}" "${MAIN_PKG}"
}

build_target linux amd64 "${BUILD_DIR}/${APP_NAME}-linux-amd64"
build_target windows amd64 "${BUILD_DIR}/${APP_NAME}-windows-amd64.exe"
build_target darwin amd64 "${BUILD_DIR}/${APP_NAME}-darwin-amd64"

print_step "Creating source archive..."
tar -czf "${BUILD_DIR}/${APP_NAME}-src.tar.gz" -C "${ROOT_DIR}" src

echo
echo "✅ Build completed successfully. Artifacts are available in '${BUILD_DIR}'."
