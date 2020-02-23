#!/bin/bash

cd $(dirname $0)

export root=$(realpath .)
export fs=$(realpath fs)
export cert=$(realpath cert)
export src=$(realpath ../src)
export ghfs=$(realpath bin)/ghfs

go build -o "$ghfs" "$src/main.go"


for file in case/*.bash; do
	bash "$file"
done;

killall "$ghfs" 2> /dev/null
