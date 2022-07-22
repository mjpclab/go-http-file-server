#!/bin/bash

cd $(dirname "$0")

export GO111MODULE=auto
source ./build.inc.sh

mkdir -p "$OUTDIR"

for build in "$@"; do
	arg=($build)
	export GOOS="${arg[0]}"
	export GOARCH="${arg[1]}"
	OS_SUFFIX="${arg[2]}"

	if [ -n "$GOLANG_VERSION" ]; then	# in docker container
		cd /usr/src/go/src/
		bash make.bash
		cd -
	fi;

	cp -r ../src/ /tmp/
	sed -i -e '/var appVer/s/"dev"/"'$VERSION'"/' /tmp/src/version/main.go

	BIN="$TMP/$MAINNAME$(go env GOEXE)"
	rm -f "$BIN"
	echo "Building: $GOOS$OS_SUFFIX $GOARCH"
	go build -ldflags "$LDFLAGS" -o "$BIN" /tmp/src/main.go

	OUT="$OUTDIR/$MAINNAME-$VERSION-$GOOS$OS_SUFFIX-$GOARCH".zip
	zip -j "$OUT" "$BIN" "$LICENSE" "$LICENSE_GO"
done
