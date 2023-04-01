#!/bin/bash

cleanup() {
	rm -f "$fs"/uploaded/[12]/*.tmp
}

source "$root"/lib.bash

"$ghfs" -l 3003 -r "$fs"/uploaded --delete /1 --delete-dir "$fs"/uploaded/2 &
sleep 0.05 # wait server ready
cleanup

file1="$fs"/uploaded/1/1.tmp
echo -n 'uploaded/1/1.tmp' > "$file1"
ls "$file1" &> /dev/null || fail "$file1 not exists"
curl_post_status -d 'name=1.tmp' 'http://127.0.0.1:3003/1?delete' > /dev/null
ls "$file1" &> /dev/null && fail "$file1 exists"

file2="$fs"/uploaded/2/2.tmp
echo -n 'uploaded/2/2.tmp' > "$file2"
ls "$file2" &> /dev/null || fail "$file2 not exists"
curl_post_status -d 'name=2.tmp' 'http://127.0.0.1:3003/2?delete' > /dev/null
ls "$file2" &> /dev/null && fail "$file2 exists"

cleanup
jobs -p | xargs kill &> /dev/null
