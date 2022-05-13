#!/bin/bash

cd $(dirname "$0")
rm -rf ../output/

builds=('linux 386' 'linux amd64' 'linux arm' 'linux arm64' 'linux riscv64' 'windows 386' 'windows amd64' 'windows arm' 'windows arm64' 'darwin amd64' 'darwin arm64')
builds=("${builds[@]}" 'freebsd 386' 'freebsd amd64' 'freebsd arm' 'freebsd arm64')
builds=("${builds[@]}" 'openbsd 386' 'openbsd amd64' 'openbsd arm' 'openbsd arm64')
builds=("${builds[@]}" 'netbsd 386' 'netbsd amd64' 'netbsd arm' 'netbsd arm64')
builds=("${builds[@]}" 'dragonfly amd64')
builds=("${builds[@]}" 'plan9 386' 'plan9 amd64' 'plan9 arm')
bash ./build.sh "${builds[@]}"
