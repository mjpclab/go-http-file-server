#!/bin/bash

source "$root"/lib.bash

"$ghfs" -l 3003 -r "$fs"/vhost1 -A --global-restrict-access example1.com example2.com &
"$ghfs" -l 3004 -r "$fs"/vhost1 -A --restrict-access @/hello@example1.com@example2.com &
sleep 0.05 # wait server ready

# --global-restrict-access

# dir list page is still allowed without referer
assert $(curl_head_status 'http://127.0.0.1:3003/hello/') '200'

# archive exposes file content, so it should be restricted as file content
assert $(curl_head_status 'http://127.0.0.1:3003/hello/?tar') '403'
assert $(curl_head_status --referer 'http://foobar.com/' 'http://127.0.0.1:3003/hello/?tar') '403'
assert $(curl_head_status --referer 'http://example1.com/' 'http://127.0.0.1:3003/hello/?tar') '200'
assert $(curl_head_status --referer 'http://example2.com/' 'http://127.0.0.1:3003/hello/?tar') '200'
assert $(curl_head_status --referer 'http://127.0.0.1:3003/hello/' 'http://127.0.0.1:3003/hello/?tar') '200'

# restricted file content should not leak out of the archive
(curl_get_body 'http://127.0.0.1:3003/hello/?tar' | tar -tf - 2> /dev/null | grep -q 'index.txt') &&
	fail "restricted content should not leak via archive"

# --restrict-access

assert $(curl_head_status 'http://127.0.0.1:3004/hello/?tar') '403'
assert $(curl_head_status --referer 'http://foobar.com/' 'http://127.0.0.1:3004/hello/?tar') '403'
assert $(curl_head_status --referer 'http://example1.com/' 'http://127.0.0.1:3004/hello/?tar') '200'
assert $(curl_head_status --referer 'http://127.0.0.1:3004/hello/' 'http://127.0.0.1:3004/hello/?tar') '200'

# path not covered by --restrict-access is unaffected
assert $(curl_head_status 'http://127.0.0.1:3004/world/?tar') '200'

jobs -p | xargs kill &> /dev/null
