#!/bin/bash

cd $(dirname "$0")
rm -rf ../output/

bash ./build.sh "$(go env GOOS) $(go env GOARCH)"
