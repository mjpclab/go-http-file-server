#!/bin/bash

source "$root"/lib.bash

"$ghfs" -l 3003 -r "$fs"/vhost1 -a :world/go:"$fs"/vhost1/go :/not/exist/name:"$fs"/vhost1/yes -a :/foo:"$fs"/vhost2/a -a :/foo/bar/baz:"$fs"/vhost2/b -E '' &
sleep 0.05 # wait server ready

file1=$(curl_get_body http://127.0.0.1:3003/file1.txt)
assert "$file1" 'vhost1/file1.txt'

hello=$(curl_get_body http://127.0.0.1:3003/hello/index.txt)
assert "$hello" 'vhost1/hello/index.txt'

go=$(curl_get_body http://127.0.0.1:3003/world/go/index.txt)
assert "$go" 'vhost1/go/index.txt'

status=$(curl_get_status http://127.0.0.1:3003/not)
assert "$status" '404'

status=$(curl_get_status http://127.0.0.1:3003/not/exist)
assert "$status" '404'

yes=$(curl_get_body http://127.0.0.1:3003/not/exist/name/index.txt)
assert "$yes" 'vhost1/yes/index.txt'

a=$(curl_get_body http://127.0.0.1:3003/foo/a1.txt)
assert "$a" 'vhost2/a/a1.txt'

status=$(curl_get_status http://127.0.0.1:3003/foo/bar)
assert "$status" '404'

baz=$(curl_get_body http://127.0.0.1:3003/foo/bar/baz/b1.txt)
assert "$baz" 'vhost2/b/b1.txt'

jobs -p | xargs kill
