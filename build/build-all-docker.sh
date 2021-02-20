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
builds=('linux 386' 'linux amd64' 'linux arm' 'linux arm64' 'windows 386' 'windows amd64' 'windows arm' 'darwin amd64' 'darwin arm64')
builds=("${builds[@]}" 'freebsd 386' 'freebsd amd64' 'freebsd arm' 'freebsd arm64')
builds=("${builds[@]}" 'openbsd 386' 'openbsd amd64' 'openbsd arm' 'openbsd arm64')
builds=("${builds[@]}" 'netbsd 386' 'netbsd amd64' 'netbsd arm' 'netbsd arm64')
builds=("${builds[@]}" 'plan9 386' 'plan9 amd64' 'plan9 arm')
buildByDocker "$gover" "${builds[@]}"

#gover=1.9
#builds=('freebsd 386 -9.x' 'freebsd amd64 -9.x' 'freebsd arm -9.x')
#buildByDocker "$gover" "${builds[@]}"

gover=1.10
builds=('windows 386 -xp-vista' 'windows amd64 -xp-vista')
#builds=("${builds[@]}" 'darwin 386 -10.8-mountain_lion' 'darwin amd64 -10.8-mountain_lion')
#builds=("${builds[@]}" 'openbsd 386 -6.0' 'openbsd amd64 -6.0' 'openbsd arm -6.0')
buildByDocker "$gover" "${builds[@]}"

#gover=1.12
#builds=('darwin 386 -10.10-yosemite' 'darwin amd64 -10.10-yosemite')
#builds=("${builds[@]}" 'freebsd 386 -10.x' 'freebsd amd64 -10.x' 'freebsd arm -9.x')
#buildByDocker "$gover" "${builds[@]}"

#gover=1.14
#builds=('darwin 386 -10.11-el_capitan' 'darwin amd64 -10.11-el_capitan')
#buildByDocker "$gover" "${builds[@]}"

#gover=1.16
#builds=('darwin amd64 -10.12-sierra')
#buildByDocker "$gover" "${builds[@]}"
