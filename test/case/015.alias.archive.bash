#!/bin/bash

cleanup() {
	rm -f "$fs"/downloaded/*.tmp
}

source "$root"/lib.bash

"$ghfs" -l 3003 -r "$fs"/vhost2 --archive / -a :/go:"$fs"/vhost1/go -a :/hello/world:"$fs"/vhost1/world -E '' &
sleep 0.05 # wait server ready
cleanup

archive="$fs"/downloaded/archive.tar.tmp
curl_get_body 'http://127.0.0.1:3003/?tar' > "$archive"
(tar -tf "$archive" | grep -q 'go/index.txt') || fail "go/index.txt not in $(basename $archive)"
(tar -tf "$archive" | grep -q 'hello/world/index.txt') || fail "hello/world/index.txt not in $(basename $archive)"

cleanup
kill %1
