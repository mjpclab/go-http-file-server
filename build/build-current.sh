#!/bin/bash

cd "$(dirname "$0")"
rm -rf ../output/

GOARCH=$(go env GOARCH)
ARCH_OPT_NAME=$(echo "GO$GOARCH" | tr 'a-z' 'A-Z')
ARCH_OPT_VALUE=$(go env "$ARCH_OPT_NAME")
if [ -n "$ARCH_OPT_VALUE" ]; then
	ARCH_OPT=",$ARCH_OPT_VALUE"
fi

bash ./build.sh "$(go env GOOS) ${GOARCH}${ARCH_OPT}"
