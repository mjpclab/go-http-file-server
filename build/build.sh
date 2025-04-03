#!/bin/bash

cd "$(dirname "$0")"

source ./build.inc.sh

mkdir -p "$OUTDIR"

for build in "$@"; do
	arg=($build)
	export GOOS="${arg[0]}"
	ARCH="${arg[1]}"	# e.g. "amd64" or "amd64,v2"
	export GOARCH=${ARCH%,*}
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
	sed -i -e '/var appArch/s/runtime.GOARCH/"'$ARCH'"/' /tmp/src/version/main.go
	mount --bind /tmp/src ../src

	TMP=$(mktemp -d)

	echo "Building: $GOOS$OS_SUFFIX $ARCH"
	go build -ldflags "$(getLdFlags)" -o "$TMP/$MAINNAME$(go env GOEXE)" ../main.go
	cp ../LICENSE "$TMP"
	cp ../src/shimgo/LICENSE_GO "$TMP"

	OUTFILE="$OUTDIR/$MAINNAME-$VERSION-$GOOS$OS_SUFFIX-$GOARCH$ARCH_OPT"
	if [ "$GOOS" == "windows" ]; then
		zip -qrj "${OUTFILE}.zip" "$TMP/"
	else
		tar --owner=0 --group=0 -zcf "${OUTFILE}.tar.gz" -C "$TMP" $(ls -A1 "$TMP")
	fi
done
