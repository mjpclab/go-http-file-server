#!/bin/bash

cleanup() {
	rm -f "$fs"/downloaded/*.tmp
}

source "$root"/lib.bash

"$ghfs" -l 3003 -r "$fs"/vhost2 --archive /a --archive-dir "$fs"/vhost2/b &
sleep 0.05 # wait server ready
cleanup

archive="$fs"/downloaded/a.tar.tmp
curl_get_body 'http://127.0.0.1:3003/a?tar' > "$archive"
(tar -tf "$archive" | grep -q a1.txt) || fail "a1.txt not in $(basename $archive)"
(tar -tf "$archive" | grep -q a2.txt) || fail "a2.txt not in $(basename $archive)"

archive="$fs"/downloaded/b.tar.tmp
curl_get_body 'http://127.0.0.1:3003/b/?tar' > "$archive"
(tar -tf "$archive" | grep -q b1.txt) || fail "b1.txt not in $(basename $archive)"
(tar -tf "$archive" | grep -q b2.txt) || fail "b2.txt not in $(basename $archive)"

cleanup
jobs -p | xargs kill
