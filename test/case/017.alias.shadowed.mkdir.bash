#!/bin/bash

cleanup() {
	rm -f "$fs"/uploaded/[12]/*.tmp
}

source "$root"/lib.bash

"$ghfs" -l 3003 -r "$fs"/uploaded/1 -a :/x.tmp:"$fs"/uploaded/2 -a:/y.tmp/z.tmp:"$fs"/uploaded/2 --delete / &
sleep 0.05 # wait server ready
cleanup

file1="$fs"/uploaded/1/x.tmp
ls "$file1" &> /dev/null && fail "$file1 exists"
curl_head_status 'http://127.0.0.1:3003/?mkdir&name=x.tmp' > /dev/null
ls "$file1" &> /dev/null && fail "$file1 exists"

file2="$fs"/uploaded/1/y.tmp
ls "$file2" &> /dev/null && fail "$file2 exists"
curl_head_status 'http://127.0.0.1:3003/?mkdir&name=x.tmp' > /dev/null
ls "$file2" &> /dev/null && fail "$file2 exists"

cleanup
jobs -p | xargs kill
