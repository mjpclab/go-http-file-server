#!/bin/bash

source "$root"/lib.bash

"$ghfs" -l 3003 -r "$fs/vhost1" --global-restrict-access example1.com example2.com &
sleep 0.05 # wait server ready

assert $(curl_head_status 'http://127.0.0.1:3003/') '200'
assert $(curl_head_status 'http://127.0.0.1:3003/hello') '200'

assert $(curl_head_status 'http://127.0.0.1:3003/hello/index.txt') '403'
assert $(curl_head_status --referer 'http://foobar.com/' 'http://127.0.0.1:3003/hello/index.txt') '403'
assert $(curl_head_status --referer 'http://example1.com/' 'http://127.0.0.1:3003/hello/index.txt') '200'
assert $(curl_head_status --referer 'http://example2.com/' 'http://127.0.0.1:3003/hello/index.txt') '200'
assert $(curl_head_status --referer 'http://127.0.0.1:3003/hello/' 'http://127.0.0.1:3003/hello/index.txt') '200'

jobs -p | xargs kill &> /dev/null
