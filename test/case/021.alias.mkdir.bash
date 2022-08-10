#!/bin/bash

cleanup() {
	rm -rf "$fs"/uploaded/[12]/*.tmp
}

source "$root"/lib.bash

"$ghfs" -l 3003 -r "$fs"/vhost1 -a :/free/space:"$fs"/uploaded/1 -a :/share/dir:"$fs"/uploaded/2 --mkdir /free/space --mkdir-dir "$fs"/uploaded/2 &
sleep 0.05 # wait server ready
cleanup

file1="$fs"/uploaded/1/foo.tmp
curl_head_status 'http://127.0.0.1:3003/free/space?mkdir&name=foo.tmp' > /dev/null
[ -d "$file1" ] || fail "$file1 should exists as directory"

file2="$fs"/uploaded/2/bar.tmp
curl_head_status 'http://127.0.0.1:3003/share/dir/?mkdir&name=bar.tmp' > /dev/null
[ -d "$file2" ] || fail "$file2 should exists as directory"

cleanup
jobs -p | xargs kill &> /dev/null
