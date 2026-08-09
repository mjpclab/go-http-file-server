#!/bin/bash

cleanup() {
	rm -f "$fs"/downloaded/*.tmp
}

source "$root"/lib.bash

# "-I go" makes the dir index lookup resolve to the "/go" alias even on empty root,
# archiving should not be affected by it
"$ghfs" -l 3003 -R --archive / -a :/go:"$fs"/vhost1/go -a :/hello/world:"$fs"/vhost1/world -I go -E '' &
sleep 0.05 # wait server ready
cleanup

archive="$fs"/downloaded/dir-index.tar.tmp
curl_get_body 'http://127.0.0.1:3003/?tar' > "$archive"
(tar -tf "$archive" | grep -q '^go/index.txt$') || fail "go/index.txt should in $(basename $archive)"
(tar -tf "$archive" | grep -q '^hello/world/index.txt$') || fail "hello/world/index.txt should in $(basename $archive)"

cleanup
jobs -p | xargs kill &> /dev/null
