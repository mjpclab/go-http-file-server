#!/bin/bash

cd "$(dirname "$0")"
source ./build.inc.version.sh

tags=(
	-t "mjpclab/ghfs:latest"
	-t "mjpclab/ghfs:$VERSION"
	-t "ghcr.io/mjpclab/ghfs:latest"
	-t "ghcr.io/mjpclab/ghfs:$VERSION"
)

docker buildx create --name ghfs-builder --driver docker-container --bootstrap
docker buildx use ghfs-builder

docker buildx build \
	"${tags[@]}" \
	-f ./build-docker-image-dockerfile \
	--platform linux/amd64,linux/386,linux/arm64,linux/arm/v7,linux/riscv64 \
	--push \
	../
