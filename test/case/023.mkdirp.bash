#!/bin/bash

cleanup() {
	rm -rf "$fs"/uploaded/[12]/*.tmp
}

source "$root"/lib.bash

"$ghfs" -l 3003 -r "$fs"/uploaded/ --mkdir /1 -E '' &
sleep 0.05 # wait server ready
cleanup

file1="$fs"/uploaded/1/foo/bar/baz.tmp
curl_post_status -d 'name=foo/bar/baz.tmp' 'http://127.0.0.1:3003/1/?mkdir' > /dev/null
[ -d "$file1" ] || fail "$file1 should exists as directory"

cleanup
jobs -p | xargs kill &> /dev/null
