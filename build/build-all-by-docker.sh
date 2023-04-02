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
        dpkg -i '"$ghfs"'/build/pkg/*.deb
      elif [ -e /etc/apk/repositories ]; then
        sed -i "s;://[^/ ]*;://mirrors.aliyun.com;" /etc/apk/repositories
        apk add bash git zip
      fi
      git config --global safe.directory "*"
      /bin/bash '"$ghfs"'/build/build.sh "$@";
      chown -R $EX_UID:$EX_GID '"$ghfs"'/output;
    ' \
    'argv_0_placeholder' \
    "$@"
}

#gover=1.15
#builds=()
#builds+=('linux 386' 'linux amd64' 'linux arm' 'linux arm64' 'windows 386' 'windows amd64' 'windows arm' 'darwin amd64' 'darwin arm64')
#builds+=('freebsd 386' 'freebsd amd64' 'freebsd arm' 'freebsd arm64')
#builds+=('openbsd 386' 'openbsd amd64' 'openbsd arm' 'openbsd arm64')
#builds+=('netbsd 386' 'netbsd amd64' 'netbsd arm' 'netbsd arm64')
#buildByDocker "$gover" "${builds[@]}"

#gover=1.15-alpine
#builds=('linux amd64 -musl' 'linux arm64 -musl')
#buildByDocker "$gover" "${builds[@]}"

#gover=1.14
#builds=('darwin 386 -10.11-el_capitan' 'darwin amd64 -10.11-el_capitan')
#buildByDocker "$gover" "${builds[@]}"

#gover=1.12
#builds=()
#builds+=('darwin 386 -10.10-yosemite' 'darwin amd64 -10.10-yosemite')
#builds+=('freebsd 386 -10.x' 'freebsd amd64 -10.x' 'freebsd arm -9.x')
#buildByDocker "$gover" "${builds[@]}"

gover=1.10
builds=()
builds+=('windows 386 -xp-vista' 'windows amd64 -xp-vista')
#builds+=('darwin 386 -10.8-mountain_lion' 'darwin amd64 -10.8-mountain_lion')
#builds+=('openbsd 386 -6.0' 'openbsd amd64 -6.0' 'openbsd arm -6.0')
buildByDocker "$gover" "${builds[@]}"

#gover=1.9
#builds=('freebsd 386 -9.x' 'freebsd amd64 -9.x' 'freebsd arm -9.x')
#buildByDocker "$gover" "${builds[@]}"
