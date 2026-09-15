#!/usr/bin/env bash
set -euo pipefail

root="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$root"

gobin="$(go env GOBIN)"
if [[ -z "$gobin" ]]; then
  gobin="$(go env GOPATH)/bin"
fi

mkdir -p "$gobin"

target="${gobin}/martini"
goos="${GOOS:-$(go env GOOS)}"
goarch="${GOARCH:-$(go env GOARCH)}"

echo "Building static martini for ${goos}/${goarch}..."
CGO_ENABLED=0 GOOS="$goos" GOARCH="$goarch" go build \
  -trimpath \
  -ldflags="-s -w" \
  -o "$target" \
  .

echo "Installed ${target}"
