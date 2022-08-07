#!/bin/bash

source "$root"/lib.bash

"$ghfs" -l 3003 -r "$fs/vhost1" -a ":/vh2:$fs/vhost2" \
	--restrict-access @/hello@example1.com@example2.com @/vh2/a@example1.com@example2.com \
	--restrict-access-dir "#$fs/vhost1/world#example1.com#example2.com" "#$fs/vhost2/b#example1.com#example2.com" \
	&
sleep 0.05 # wait server ready

assert $(curl_head_status 'http://127.0.0.1:3003/') '200'
assert $(curl_head_status 'http://127.0.0.1:3003/hello') '200'
assert $(curl_head_status 'http://127.0.0.1:3003/world') '200'
assert $(curl_head_status 'http://127.0.0.1:3003/vh2') '200'
assert $(curl_head_status 'http://127.0.0.1:3003/vh2/a/') '200'
assert $(curl_head_status 'http://127.0.0.1:3003/vh2/b/') '200'

assert $(curl_head_status 'http://127.0.0.1:3003/hello/index.txt') '403'
assert $(curl_head_status --referer 'http://foobar.com/' 'http://127.0.0.1:3003/hello/index.txt') '403'
assert $(curl_head_status --referer 'http://example1.com/' 'http://127.0.0.1:3003/hello/index.txt') '200'
assert $(curl_head_status --referer 'http://example2.com/' 'http://127.0.0.1:3003/hello/index.txt') '200'
assert $(curl_head_status --referer 'http://127.0.0.1:3003/hello/' 'http://127.0.0.1:3003/hello/index.txt') '200'

assert $(curl_head_status 'http://127.0.0.1:3003/world/index.txt') '403'
assert $(curl_head_status --referer 'http://foobar.com/' 'http://127.0.0.1:3003/world/index.txt') '403'
assert $(curl_head_status --referer 'http://example1.com/' 'http://127.0.0.1:3003/world/index.txt') '200'
assert $(curl_head_status --referer 'http://example2.com/' 'http://127.0.0.1:3003/world/index.txt') '200'
assert $(curl_head_status --referer 'http://127.0.0.1:3003/world/' 'http://127.0.0.1:3003/world/index.txt') '200'

assert $(curl_head_status 'http://127.0.0.1:3003/go/index.txt') '200'

assert $(curl_head_status 'http://127.0.0.1:3003/vh2/a/a1.txt') '403'
assert $(curl_head_status --referer 'http://foobar.com/' 'http://127.0.0.1:3003/vh2/a/a1.txt') '403'
assert $(curl_head_status --referer 'http://example1.com/' 'http://127.0.0.1:3003/vh2/a/a1.txt') '200'
assert $(curl_head_status --referer 'http://example2.com/' 'http://127.0.0.1:3003/vh2/a/a1.txt') '200'
assert $(curl_head_status --referer 'http://127.0.0.1:3003/vh2/a/' 'http://127.0.0.1:3003/vh2/a/a1.txt') '200'

assert $(curl_head_status 'http://127.0.0.1:3003/vh2/b/b1.txt') '403'
assert $(curl_head_status --referer 'http://foobar.com/' 'http://127.0.0.1:3003/vh2/b/b1.txt') '403'
assert $(curl_head_status --referer 'http://example1.com/' 'http://127.0.0.1:3003/vh2/b/b1.txt') '200'
assert $(curl_head_status --referer 'http://example2.com/' 'http://127.0.0.1:3003/vh2/b/b1.txt') '200'
assert $(curl_head_status --referer 'http://127.0.0.1:3003/vh2/b/' 'http://127.0.0.1:3003/vh2/b/b1.txt') '200'

assert $(curl_head_status 'http://127.0.0.1:3003/vh2/file1.txt') '200'

jobs -p | xargs kill &> /dev/null
