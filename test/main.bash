#!/bin/bash

cd $(dirname $0)

for cmd in realpath curl grep sed xargs jq; do
	type "$cmd" &> /dev/null
	if [ $? -ne 0 ]; then
		echo "command '$cmd' not found" >&2
		exit 1
	fi
done

export root=$(realpath .)
export fs=$(realpath fs)
export cert=$(realpath cert)
export ghfs=$(realpath bin)/ghfs
export GHFS_QUIET=1

go build -o "$ghfs" "../main.go"

pattern="$1"
if [ -z "$pattern" ]; then
	pattern='*'
fi

for file in case/$pattern.bash; do
	bash "$file"
done

killall "$ghfs" 2> /dev/null
