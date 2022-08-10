#!/bin/bash

cleanup() {
	rm -rf "$fs"/uploaded/[12]/*.tmp
}

source "$root"/lib.bash

"$ghfs" -l 3003 -r "$fs"/uploaded/ --mkdir /1 &
sleep 0.05 # wait server ready
cleanup

file1="$fs"/uploaded/1/foo.tmp
curl_head_status 'http://127.0.0.1:3003/1/?mkdir&name=foo.tmp' > /dev/null
[ -d "$file1" ] || fail "$file1 should exists as directory"

cleanup
jobs -p | xargs kill &> /dev/null
