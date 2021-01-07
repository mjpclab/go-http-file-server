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
      apt-get install -yq --force-yes git zip;
      /bin/bash /mnt/build/build.sh "$@";
      chown -R $EX_UID:$EX_GID /mnt/output
    ' \
    'container-script' \
    "$@"
}

gover=1.2
builds=('windows 386 -2000')
#builds=("${builds[@]}" 'freebsd 386 -8' 'freebsd amd64 -8')
buildByDocker "$gover" "${builds[@]}"
