#!/bin/bash

cleanup() {
	rm -f "$fs"/downloaded/*.tmp
}

source "$root"/lib.bash

"$ghfs" -l 3003 -r "$fs"/vhost2 --archive /a --archive-dir "$fs"/vhost2/b --archive /c --auth /c/sub -a :/hello:"$fs"/vhost1/hello -a :/world:"$fs"/vhost1/world --archive-user :/hello:alice --archive-dir-user :"$fs"/vhost1/world:bob --user alice:AliceSecret bob:BobSecret &
sleep 0.05 # wait server ready
cleanup

archive="$fs"/downloaded/a.tar.tmp
curl_get_body 'http://127.0.0.1:3003/a?tar' > "$archive"
(tar -tf "$archive" | grep -q '^a1.txt$') || fail "a1.txt should in $(basename $archive)"
(tar -tf "$archive" | grep -q '^a2.txt$') || fail "a2.txt should in $(basename $archive)"

(curl_get_header 'http://127.0.0.1:3003/a?tar=saved.tar' | grep -a -i -F 'content-disposition' | grep -q -F 'saved.tar') ||
	fail "archive specified filename is not exists"

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

assert $(curl_head_status 'http://127.0.0.1:3003/hello/?tar') '400'
assert $(curl_head_status 'http://eve:EveSecret@127.0.0.1:3003/hello/?tar') '400'
assert $(curl_head_status 'http://alice:WrongPass@127.0.0.1:3003/hello/?tar') '400'
assert $(curl_head_status 'http://alice:AliceSecret@127.0.0.1:3003/hello/?tar') '200'
archive="$fs"/downloaded/hello.tar.tmp
curl_get_body 'http://alice:AliceSecret@127.0.0.1:3003/hello/?tar' > "$archive"
(tar -tf "$archive" | grep -q '^index.txt$') || fail "index.txt should in $(basename $archive)"

assert $(curl_head_status 'http://127.0.0.1:3003/world/?tar') '400'
assert $(curl_head_status 'http://eve:EveSecret@127.0.0.1:3003/world/?tar') '400'
assert $(curl_head_status 'http://bob:WrongPass@127.0.0.1:3003/world/?tar') '400'
assert $(curl_head_status 'http://bob:BobSecret@127.0.0.1:3003/world/?tar') '200'
archive="$fs"/downloaded/world.tar.tmp
curl_get_body 'http://bob:BobSecret@127.0.0.1:3003/world/?tar' > "$archive"
(tar -tf "$archive" | grep -q '^index.txt$') || fail "index.txt should in $(basename $archive)"

cleanup
jobs -p | xargs kill &> /dev/null
