#!/bin/bash

cd $(dirname "$0")
rm -rf ../output/

builds=('linux 386' 'linux amd64' 'linux arm' 'linux arm64' 'windows 386 .exe' 'windows amd64 .exe' 'darwin 386' 'darwin amd64')
bash ./build.sh "${builds[@]}"
