#!/bin/bash

prefix=$(dirname "$0")/../
absPrefix=$(realpath "$prefix")

rm -rf "$prefix/output/"

buildByDocker() {
  gover="$1"
  shift
  docker pull golang:"$gover"

  docker run \
    --rm \
    -v "$absPrefix":/mnt \
    -e EX_UID="$(id -u)" \
    -e EX_GID="$(id -g)" \
    golang:"$gover" \
    /bin/bash -c '
      sed -i -e "s;://[^/ ]*;://mirrors.aliyun.com;" /etc/apt/sources.list;
      apt-get update;
      apt-get install -yq git zip;
      /bin/bash /mnt/build/build.sh "$@";
      chown -R $EX_UID:$EX_GID /mnt/output
    ' \
    'container-script' \
    "$@"
}

gover=latest
builds=('linux 386' 'linux amd64' 'linux arm' 'linux arm64' 'windows 386' 'windows amd64' 'windows arm' 'windows arm64' 'darwin amd64' 'darwin arm64')
builds=("${builds[@]}" 'freebsd 386' 'freebsd amd64' 'freebsd arm' 'freebsd arm64')
builds=("${builds[@]}" 'openbsd 386' 'openbsd amd64' 'openbsd arm' 'openbsd arm64')
builds=("${builds[@]}" 'netbsd 386' 'netbsd amd64' 'netbsd arm' 'netbsd arm64')
builds=("${builds[@]}" 'dragonfly amd64')
builds=("${builds[@]}" 'plan9 386' 'plan9 amd64' 'plan9 arm')
buildByDocker "$gover" "${builds[@]}"

#gover=1.16
#builds=('darwin amd64 -10.12-sierra')
#buildByDocker "$gover" "${builds[@]}"
