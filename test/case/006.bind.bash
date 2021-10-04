#!/bin/bash

source "$root"/lib.bash

"$ghfs" -l 3003 -r "$fs"/vhost1 -b :/foo:"$fs"/vhost2/a -b :/foo/bar/baz:"$fs"/vhost2/b -E '' &
sleep 0.05 # wait server ready

a=$(curl_get_body http://127.0.0.1:3003/foo/a1.txt)
assert "$a" 'vhost2/a/a1.txt'

a=$(curl_get_body http://127.0.0.1:3003/Foo/a1.txt)
assert "$a" 'vhost2/a/a1.txt'

status=$(curl_get_status http://127.0.0.1:3003/foo/bar)
assert "$status" '404'

status=$(curl_get_status http://127.0.0.1:3003/Foo/Bar)
assert "$status" '404'

baz=$(curl_get_body http://127.0.0.1:3003/foo/bar/baz/b1.txt)
assert "$baz" 'vhost2/b/b1.txt'

baz=$(curl_get_body http://127.0.0.1:3003/foo/bar/Baz/b1.txt)
assert "$baz" 'vhost2/b/b1.txt'

jobs -p | xargs kill
