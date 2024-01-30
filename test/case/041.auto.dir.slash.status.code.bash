#!/bin/bash

source "$root"/lib.bash

"$ghfs" -l 3003 --prefix foo/bar -r "$fs/vhost1" -a ":/lorem/ipsum:$fs/vhost2" --dir-index index.txt --auto-dir-slash 302 -E '' &
sleep 0.05 # wait server ready

assert $(curl_head_status 'http://127.0.0.1:3003/foo/bar') '302'
(curl_get_header http://127.0.0.1:3003/foo/bar | grep -q -i 'location:\s*/foo/bar/') ||
	fail "incorrect redirect location"
assert $(curl_head_status 'http://127.0.0.1:3003/foo/bar/') '200'

assert $(curl_head_status 'http://127.0.0.1:3003/foo/bar/hello?sort=T') '302'
(curl_get_header 'http://127.0.0.1:3003/foo/bar/hello?sort=T' | grep -q -i 'location:\s*/foo/bar/hello/?sort=T') ||
	fail "incorrect redirect location"
assert $(curl_head_status 'http://127.0.0.1:3003/foo/bar/hello/') '200'

assert $(curl_head_status 'http://127.0.0.1:3003/foo/bar/lorem') '302'
(curl_get_header http://127.0.0.1:3003/foo/bar/lorem | grep -q -i 'location:\s*/foo/bar/lorem/') ||
	fail "incorrect redirect location"
assert $(curl_head_status 'http://127.0.0.1:3003/foo/bar/lorem/') '404'

assert $(curl_head_status 'http://127.0.0.1:3003/foo/bar/lorem/ipsum') '302'
(curl_get_header http://127.0.0.1:3003/foo/bar/lorem/ipsum | grep -q -i 'location:\s*/foo/bar/lorem/ipsum/') ||
	fail "incorrect redirect location"
assert $(curl_head_status 'http://127.0.0.1:3003/foo/bar/lorem/ipsum/') '200'

assert $(curl_head_status 'http://127.0.0.1:3003/foo/bar/lorem/ipsum/a') '302'
(curl_get_header http://127.0.0.1:3003/foo/bar/lorem/ipsum/a | grep -q -i 'location:\s*/foo/bar/lorem/ipsum/a/') ||
	fail "incorrect redirect location"
assert $(curl_head_status 'http://127.0.0.1:3003/foo/bar/lorem/ipsum/a/') '200'

assert $(curl_head_status 'http://127.0.0.1:3003/foo/bar/lorem/ipsum/a/not-exist') '404'


assert $(curl_head_status 'http://127.0.0.1:3003/foo/bar/file1.txt/') '302'
(curl_get_header http://127.0.0.1:3003/foo/bar/file1.txt/ | grep -q -i 'location:\s*/foo/bar/file1.txt\s*$') ||
	fail "incorrect redirect location"
assert $(curl_head_status 'http://127.0.0.1:3003/foo/bar/file1.txt') '200'

assert $(curl_head_status 'http://127.0.0.1:3003/foo/bar/lorem/ipsum/file1.txt/') '302'
(curl_get_header http://127.0.0.1:3003/foo/bar/lorem/ipsum/file1.txt/ | grep -q -i 'location:\s*/foo/bar/lorem/ipsum/file1.txt\s*$') ||
	fail "incorrect redirect location"
assert $(curl_head_status 'http://127.0.0.1:3003/foo/bar/lorem/ipsum/file1.txt') '200'

jobs -p | xargs kill &> /dev/null
