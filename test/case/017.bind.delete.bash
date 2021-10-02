#!/bin/bash

cleanup() {
	rm -f "$fs"/uploaded/[12]/*.tmp
}

source "$root"/lib.bash

"$ghfs" -l 3003 -r "$fs"/vhost1 -b :/free/space:"$fs"/uploaded/1 -b :/share/dir:"$fs"/uploaded/2 --delete /free/space --delete-dir "$fs"/uploaded/2 &
sleep 0.05 # wait server ready
cleanup

file1="$fs"/uploaded/1/1.tmp
echo -n 'uploaded/1/1.tmp' > "$file1"
[ -e "$file1" ] || fail "$file1 should exists"
curl_head_status 'http://127.0.0.1:3003/free/space?delete&name=1.tmp' > /dev/null
[ -e "$file1" ] && fail "$file1 should not exists"

file2="$fs"/uploaded/2/2.tmp
echo -n 'uploaded/2/2.tmp' > "$file2"
[ -e "$file2" ] || fail "$file2 should exists"
curl_head_status 'http://127.0.0.1:3003/SHARE/dir?delete&name=2.tmp' > /dev/null
[ -e "$file2" ] && fail "$file2 should not exists"

cleanup
jobs -p | xargs kill
