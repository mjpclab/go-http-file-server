#!/bin/bash

source "$root"/lib.bash

"$ghfs" -l 3003 -r "$fs"/vhost1 &
sleep 0.05 # wait server ready

file1=$(curl_get_body http://127.0.0.1:3003/file1.txt)
assert "$file1" 'vhost1/file1.txt'

(curl_get_body http://127.0.0.1:3003/hello | grep -q './hello/index.txt') ||
	fail "resource /hello does not contains './hello/index.txt'"

(curl_get_body http://127.0.0.1:3003/hello/ | grep -q './index.txt') ||
	fail "resource /hello/ does not contains './index.txt'"

(curl_get_body 'http://127.0.0.1:3003/' | grep -q -F 'escape-escaped%252fslash') ||
	fail "link URL escape incorrect"
(curl_get_body 'http://127.0.0.1:3003/' | grep -q -F 'escape%23sharp') ||
	fail "link URL escape incorrect"
(curl_get_body 'http://127.0.0.1:3003/' | grep -q -F 'escape%25percent') ||
	fail "link URL escape incorrect"

assert $(curl_head_status 'http://127.0.0.1:3003/escape-escaped%252fslash/') '200'
assert $(curl_head_status 'http://127.0.0.1:3003/escape%23sharp') '200'
assert $(curl_head_status 'http://127.0.0.1:3003/escape%25percent') '200'

kill %1
