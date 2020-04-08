#!/bin/bash

cleanup() {
	rm -rf "$fs"/uploaded/[12]/*.tmp
}

source "$root"/lib.bash

"$ghfs" -l 3003 -r "$fs"/uploaded/ --mkdir /1 &
sleep 0.05 # wait server ready
cleanup

file="$fs"/uploaded/1/foo.tmp
ls -d "$file1" &> /dev/null && fail "$file1 exists"
curl_head_status 'http://127.0.0.1:3003/1/?mkdir&name=foo.tmp' > /dev/null
ls -d "$file1"/ &> /dev/null || fail "$file1 not exists"

cleanup
kill %1
