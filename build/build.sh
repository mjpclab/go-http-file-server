#!/bin/bash

cd $(dirname "$0")

export GO111MODULE=auto
source ./build.inc.sh

mkdir -p "$OUTDIR"

for build in "$@"; do
	arg=($build)
	export GOOS="${arg[0]}"
	ARCH="${arg[1]}"	# e.g. "amd64" or "amd64,v2"
	GOARCH=${ARCH%,*}
	if [ "$ARCH" != "$GOARCH" ]; then
		# e.g. "GOAMD64=v2"
		ARCH_OPT="${ARCH#*,}"
		declare -x $(echo GO$GOARCH | tr 'a-z' 'A-Z')="$ARCH_OPT"
	else
		ARCH_OPT=''
		unset $(echo "GO$GOARCH" | tr 'a-z' 'A-Z')
	fi
	OS_SUFFIX="${arg[2]}"

	BIN="$TMP/$MAINNAME$(go env GOEXE)"
	rm -f "$BIN"
	echo "Building: $GOOS$OS_SUFFIX $ARCH"
	go build -ldflags "$(getLdFlags)" -o "$BIN" ../src/main.go

	OUT="$OUTDIR/$MAINNAME-$VERSION-$GOOS$OS_SUFFIX-$GOARCH$ARCH_OPT".zip
	zip -j "$OUT" "$BIN" "$LICENSE"
done
