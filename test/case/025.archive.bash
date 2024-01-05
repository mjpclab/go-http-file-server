#!/bin/bash

cleanup() {
	rm -f "$fs"/downloaded/*.tmp
}

source "$root"/lib.bash

"$ghfs" -l 3003 -r "$fs"/vhost2 --archive /a --archive-dir "$fs"/vhost2/b --archive /c --auth /c/sub &
sleep 0.05 # wait server ready
cleanup

archive="$fs"/downloaded/a.tar.tmp
curl_get_body 'http://127.0.0.1:3003/a?tar' > "$archive"
(tar -tf "$archive" | grep -q '^a1.txt$') || fail "a1.txt should in $(basename $archive)"
(tar -tf "$archive" | grep -q '^a2.txt$') || fail "a2.txt should in $(basename $archive)"

archive="$fs"/downloaded/a-part.tar.tmp
curl_get_body 'http://127.0.0.1:3003/a?tar&name=a1.txt' > "$archive"
(tar -tf "$archive" | grep -q '^a1.txt$') || fail "a1.txt should in $(basename $archive)"
(tar -tf "$archive" | grep -q 'a2.txt') && fail "a2.txt should not in $(basename $archive)"

archive="$fs"/downloaded/b.tar.tmp
curl_get_body 'http://127.0.0.1:3003/b/?tar' > "$archive"
(tar -tf "$archive" | grep -q '^b1.txt$') || fail "b1.txt should in $(basename $archive)"
(tar -tf "$archive" | grep -q '^b2.txt$') || fail "b2.txt should in $(basename $archive)"

archive="$fs"/downloaded/c.tar.tmp
curl_get_body 'http://127.0.0.1:3003/c/?tar' > "$archive"
(tar -tf "$archive" | grep -q '^c1.txt$') || fail "c1.txt should in $(basename $archive)"
(tar -tf "$archive" | grep -q '^c2.txt$') || fail "c2.txt should in $(basename $archive)"
(tar -tf "$archive" | grep -q '^sub/sub1.txt$') && fail "sub/sub1.txt should not in $(basename $archive)"
(tar -tf "$archive" | grep -q '^sub/sub2.txt$') && fail "sub/sub1.txt should not in $(basename $archive)"

cleanup
jobs -p | xargs kill &> /dev/null
