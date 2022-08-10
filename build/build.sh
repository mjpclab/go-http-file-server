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

	if [ -n "$GOLANG_VERSION" ]; then	# in docker container
		cd /usr/src/go/src/
		bash make.bash
		cd -
	fi;

	cp -r ../src/ /tmp/
	sed -i -e '/var appVer/s/"dev"/"'$VERSION'"/' /tmp/src/version/main.go
	sed -i -e '/var appArch/s/"runtime.GOARCH"/"'$ARCH'"/' /tmp/src/version/main.go

	BIN="$TMP/$MAINNAME$(go env GOEXE)"
	rm -f "$BIN"
	echo "Building: $GOOS$OS_SUFFIX $ARCH"
	go build -ldflags "$LDFLAGS" -o "$BIN" /tmp/src/main.go

	OUT="$OUTDIR/$MAINNAME-$VERSION-$GOOS$OS_SUFFIX-$GOARCH$ARCH_OPT".zip
	zip -j "$OUT" "$BIN" "$LICENSE" "$LICENSE_GO"
done
