#!/bin/bash

cd "$(dirname "$0")"
source ./build.inc.version.sh

TAG_PREFIX=${TAG_PREFIX:-mjpclab/ghfs}

docker build -t "$TAG_PREFIX:latest" -t "$TAG_PREFIX:$VERSION" -f ./build-docker-image-dockerfile ../
