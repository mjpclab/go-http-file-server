#!/bin/bash

cd $(dirname "$0")

TMP='/tmp'
LDFLAGS='-s -w'
OUTDIR='../output'
MAINNAME='ghfs'
VERSION="$(git describe --abbrev=0 --tags 2> /dev/null || git rev-parse --abbrev-ref HEAD 2> /dev/null)"
LICENSE='../LICENSE'
builds=('linux 386' 'linux amd64' 'linux arm' 'linux arm64' 'windows 386 .exe' 'windows amd64 .exe' 'darwin 386' 'darwin amd64')

mkdir -p "$OUTDIR"
rm -rf "$OUTDIR"/*

for build in "${builds[@]}"; do
	arg=($build)
	export GOOS="${arg[0]}"
	export GOARCH="${arg[1]}"
	SUFFIX="${arg[2]}"

	BIN="$TMP"/"$MAINNAME""$SUFFIX"
	rm -f "$BIN"
	go build -ldflags "$LDFLAGS" -o "$BIN" ../src/main.go

	OUT="$OUTDIR"/"$MAINNAME"-"$VERSION"-"$GOOS"-"$GOARCH".zip
	zip -j "$OUT" "$BIN" "$LICENSE"
done
