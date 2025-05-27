#!/usr/bin/env bash
# require upx
set -eu

: "${WORK_DIR:=./}"
: "${CHARS_FILE:=${1:-glyph.txt}}"

build() {
  local GOOS="${1:-linux}" GOARCH="${2:-amd64}" bin

  bin="homo${CHARS_FILE%.txt}-$GOOS-$GOARCH"
  [ "$GOOS" = windows ] && bin+=.exe

  printf 'Build:\t%-10s%-7s' "$GOOS" "$GOARCH"

  CGO_ENABLED=0 GOARCH="$GOARCH" GOOS="$GOOS" \
  GOFLAGS="-buildvcs=false -trimpath" \
    go build -ldflags="-s -w " -o "./build/$bin" "$WORK_DIR"/...

  echo "./build/$bin"
}

mkdir -p ./build
rm -f homoglyph.go homoglitch.go
go mod tidy
go run generate.go "$CHARS_FILE"

if [ -z "${3:-}" ]; then
  build darwin amd64
  build darwin arm64
  build linux 386
  build linux amd64
  build linux arm
  build linux arm64
  build windows 386
  build windows amd64
  build windows arm64
else
  build "${@:3}"
fi

rm "homo${CHARS_FILE%.txt}.go"
