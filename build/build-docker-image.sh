#!/bin/bash

cd "$(dirname "$0")"
source ./build.inc.version.sh

TAG_PREFIX=${TAG_PREFIX:-mjpclab/ghfs}

docker buildx create --name ghfs-builder --driver docker-container --bootstrap --driver-opt env.https_proxy=
docker buildx use ghfs-builder

docker buildx build \
	-t "$TAG_PREFIX:latest" \
	-t "$TAG_PREFIX:$VERSION" \
	-f ./build-docker-image-dockerfile \
	--platform linux/amd64,linux/386,linux/arm64,linux/arm/v7,linux/riscv64 \
	--push \
	--build-arg GOOS=linux \
	--build-arg GOARCH=amd64 \
	../
