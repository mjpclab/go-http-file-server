#!/bin/bash

cleanup() {
	rm -f "$fs"/uploaded/[12]/*.tmp
}

source "$root"/lib.bash

"$ghfs" -l 3003 -r "$fs"/uploaded/1 -a :/x.tmp:"$fs"/uploaded/2 -a:/y.tmp/z.tmp:"$fs"/uploaded/2 --delete / &
sleep 0.05 # wait server ready
cleanup

file1="$fs"/uploaded/1/x.tmp
echo -n 'uploaded/1/1.tmp' > "$file1"
[ -e "$file1" ] || fail "$file1 should exists"
curl_head_status 'http://127.0.0.1:3003/?delete&name=x.tmp' > /dev/null
[ -e "$file1" ] || fail "$file1 should exists"

file2="$fs"/uploaded/1/y.tmp
echo -n 'uploaded/1/y.tmp' > "$file2"
[ -e "$file2" ] || fail "$file2 should exists"
curl_head_status 'http://127.0.0.1:3003/?delete&name=y.tmp' > /dev/null
[ -e "$file2" ] || fail "$file2 should exists"

cleanup
jobs -p | xargs kill
