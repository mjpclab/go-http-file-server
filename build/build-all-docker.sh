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
      sed -i -e "s;://[^/]*/;://mirrors.aliyun.com/;" /etc/apt/sources.list;
      apt-get update && apt-get install -yq git zip;
      /bin/bash /mnt/build/build.sh "$@";
      chown -R $EX_UID:$EX_GID /mnt/output
    ' \
    'container-script' \
    "$@"
}

gover=latest
builds=('linux 386' 'linux amd64' 'linux arm' 'linux arm64' 'windows 386' 'windows amd64' 'darwin 386' 'darwin amd64')
buildByDocker "$gover" "${builds[@]}"

gover=1.10
builds=('windows 386 -xp-vista' 'windows amd64 -xp-vista')
buildByDocker "$gover" "${builds[@]}"
