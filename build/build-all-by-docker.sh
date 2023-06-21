#!/bin/bash

cd "$(dirname "$0")"

# init variable `builds`
source ./build-all.inc.sh

prefix=$(realpath ../)
ghfs=/go/src/mjpclab.dev/ghfs

rm -rf "$prefix/output/"

buildByDocker() {
  local tag="$1"
  shift
  docker pull golang:"$tag"

  docker run \
    --rm \
    -v "$prefix":"$ghfs" \
    -e EX_UID="$(id -u)" \
    -e EX_GID="$(id -g)" \
    golang:"$tag" \
    /bin/sh -c '
      if [ -e /etc/apt/sources.list ]; then
        sed -i -e "s;://[^/ ]*;://mirrors.aliyun.com;" /etc/apt/sources.list;
        apt-get update;
        apt-get install -yq git zip;
      elif [ -e /etc/apt/sources.list.d/debian.sources ]; then
        sed -i -e "s;://[^/ ]*;://mirrors.aliyun.com;" /etc/apt/sources.list.d/debian.sources;
        apt-get update;
        apt-get install -yq git zip;
      elif [ -e /etc/apk/repositories ]; then
        sed -i "s;://[^/ ]*;://mirrors.aliyun.com;" /etc/apk/repositories;
        apk add bash git zip;
      fi
      git config --global safe.directory "*"
      /bin/bash '"$ghfs"'/build/build.sh "$@";
      chown -R $EX_UID:$EX_GID '"$ghfs"'/output;
    ' \
    'argv_0_placeholder' \
    "$@"
}

gover=latest
buildByDocker "$gover" "${builds[@]}"

gover=alpine
builds=('linux 386 -musl' 'linux amd64 -musl' 'linux amd64,v2 -musl' 'linux amd64,v3 -musl' 'linux arm64 -musl')
buildByDocker "$gover" "${builds[@]}"

gover=1.20
builds=()
builds+=('windows 386 -7-8' 'windows amd64 -7-8')
#builds+=('windows amd64,v2 -7-8' 'windows amd64,v3 -7-8')
#builds+=('darwin amd64 -10.13-high-sierra')
buildByDocker "$gover" "${builds[@]}"

#gover=1.16
#builds=('darwin amd64 -10.12-sierra')
#buildByDocker "$gover" "${builds[@]}"
