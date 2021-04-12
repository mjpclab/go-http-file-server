#!/bin/bash

cleanup() {
	rm -f "$fs"/downloaded/*.tmp
}

source "$root"/lib.bash

"$ghfs" -l 3003 -R --archive / -a :/go:"$fs"/vhost1/go -a :/hello/world:"$fs"/vhost1/world -E '' &
sleep 0.05 # wait server ready
cleanup

archive="$fs"/downloaded/archive.tar.tmp
curl_get_body 'http://127.0.0.1:3003/?tar' > "$archive"
(tar -tf "$archive" | grep -q '^go/index.txt$') || fail "go/index.txt should in $(basename $archive)"
(tar -tf "$archive" | grep -q '^hello/world/index.txt$') || fail "hello/world/index.txt should in $(basename $archive)"
(tar -tf "$archive" | grep -q -E -v '^(go|hello)/') && fail "unexpected file in $(basename $archive)"

archive="$fs"/downloaded/archive-partial.tar.tmp
curl_get_body 'http://127.0.0.1:3003/?tar&name=go' > "$archive"
(tar -tf "$archive" | grep -q '^go/index.txt$') || fail "go/index.txt should in $(basename $archive)"
(tar -tf "$archive" | grep -q -E -v '^go/') && fail "unexpected file in $(basename $archive)"

cleanup
jobs -p | xargs kill
