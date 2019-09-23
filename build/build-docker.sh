#!/bin/bash

prefix=$(dirname "$0")/../
absPrefix=$(realpath "$prefix")

docker run \
  --rm \
  -v "$absPrefix":/mnt \
  -e EX_UID="$(id -u)" \
  -e EX_GID="$(id -g)" \
  golang:1.10 \
  /bin/bash -c '
		sed -i -e "s;://[^/]*/;://mirrors.aliyun.com/;" /etc/apt/sources.list
		apt-get update && apt-get install -yq git zip;
		/bin/bash /mnt/build/build.sh;
		chown -R $EX_UID:$EX_GID /mnt/output
	'
