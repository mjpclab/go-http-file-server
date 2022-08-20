#!/bin/bash

cd $(dirname "$0")

# init variable `builds`
source ./build-all.inc.sh

prefix=$(realpath ../)
ghfs=/go/src/mjpclab.dev/ghfs

rm -rf "$prefix/output/"

buildByDocker() {
  gover="$1"
  shift
  docker pull golang:"$gover"

  docker run \
    --rm \
    -v "$prefix":"$ghfs" \
    -e EX_UID="$(id -u)" \
    -e EX_GID="$(id -g)" \
    golang:"$gover" \
    /bin/bash -c '
      sed -i -e "s;://[^/ ]*;://mirrors.aliyun.com;" /etc/apt/sources.list;
      apt-get update;
      apt-get install -yq --force-yes git zip;
      /bin/bash '"$ghfs"'/build/build.sh "$@";
      chown -R $EX_UID:$EX_GID '"$ghfs"'/output;
    ' \
    'argv_0_placeholder' \
    "$@"
}

gover=1.2
builds=('windows 386 -2000')
#builds=("${builds[@]}" 'freebsd 386 -8' 'freebsd amd64 -8')
buildByDocker "$gover" "${builds[@]}"
