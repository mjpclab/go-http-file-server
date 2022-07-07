#!/bin/bash

source "$root"/lib.bash

"$ghfs" -l 3003 -r "$fs"/vhost1 --prefix /foo/bar --base /lorem/ipsum &
sleep 0.05 # wait server ready

file1=$(curl_get_body http://127.0.0.1:3003/foo/bar/file1.txt)
assert "$file1" 'vhost1/file1.txt'

file1=$(curl_get_body http://127.0.0.1:3003/lorem/ipsum/file1.txt)
assert "$file1" 'vhost1/file1.txt'

file1=$(curl_get_body http://127.0.0.1:3003/LOREM/Ipsum/file1.txt)
assert "$file1" 'vhost1/file1.txt'

assert $(curl_head_status 'http://127.0.0.1:3003') '404'
assert $(curl_head_status 'http://127.0.0.1:3003/') '404'

assert $(curl_head_status 'http://127.0.0.1:3003/foo') '404'
assert $(curl_head_status 'http://127.0.0.1:3003/foo/bar') '200'
assert $(curl_head_status 'http://127.0.0.1:3003/foo/bar/') '200'
assert $(curl_head_status 'http://127.0.0.1:3003/Foo/bar') '404'
assert $(curl_head_status 'http://127.0.0.1:3003/Foo/bar/') '404'

assert $(curl_head_status 'http://127.0.0.1:3003/lorem') '404'
assert $(curl_head_status 'http://127.0.0.1:3003/lorem/ipsum') '200'
assert $(curl_head_status 'http://127.0.0.1:3003/Lorem/Ipsum') '200'
assert $(curl_head_status 'http://127.0.0.1:3003/lorem/ipsum/') '200'
assert $(curl_head_status 'http://127.0.0.1:3003/Lorem/Ipsum/') '200'

jobs -p | xargs kill
