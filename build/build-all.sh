#!/bin/bash

cd $(dirname "$0")
rm -rf ../output/

# init variable `builds`
source ./build-all.inc.sh

bash ./build.sh "${builds[@]}"
