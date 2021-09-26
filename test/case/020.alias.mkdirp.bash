#!/bin/bash

cleanup() {
	rm -rf "$fs"/uploaded/[12]/*.tmp
}

source "$root"/lib.bash

"$ghfs" -l 3003 -r "$fs"/uploaded/ --mkdir /1 -a :/1/my1.tmp:"$fs"/uploaded/2 -a :/1/my2.tmp/backup:"$fs"/uploaded/2 &> /dev/null &
sleep 0.05 # wait server ready
cleanup

file1="$fs"/uploaded/1/my1.tmp
curl_head_status 'http://127.0.0.1:3003/1/?mkdir&name=my1.tmp' > /dev/null
[ -e "$file1" ] && fail "$file1 should not exists"

file2="$fs"/uploaded/2/my2.tmp
curl_head_status 'http://127.0.0.1:3003/1/?mkdir&name=my2.tmp/test' > /dev/null
[ -e "$file2" ] && fail "$file2 should not exists"

cleanup
jobs -p | xargs kill
