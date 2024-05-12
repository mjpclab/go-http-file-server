#!/bin/bash

source "$root"/lib.bash

"$ghfs" -l 3003 -r "$fs"/vhost1 -a :/x/y/z:"$fs"/vhost2 --user foo:123 bar:456 --index /hello --index-dir "$fs"/vhost1/world --index-user :/x/y/z/a:foo --index-dir-user :"$fs"/vhost2/b:bar -E '' &
sleep 0.05 # wait server ready

# --index

cnt=$(curl_get_body 'http://127.0.0.1:3003/?json' | jq '.subItems | length')
assert "$cnt" '0'

cnt=$(curl_get_body 'http://127.0.0.1:3003/go/?json' | jq '.subItems | length')
assert "$cnt" '0'

cnt=$(curl_get_body 'http://127.0.0.1:3003/hello/?json' | jq '.subItems | length')
[ "$cnt" == "0" ] && fail "subItems should not be 0"

# --index-dir

cnt=$(curl_get_body 'http://127.0.0.1:3003/world/?json' | jq '.subItems | length')
[ "$cnt" == "0" ] && fail "subItems should not be 0"

# --index-user

cnt=$(curl_get_body 'http://127.0.0.1:3003/x/y/z/a/?json' | jq '.subItems | length')
assert "$cnt" '0'

cnt=$(curl_get_body 'http://baz:789@127.0.0.1:3003/x/y/z/a/?json' | jq '.subItems | length')
assert "$cnt" '0'

cnt=$(curl_get_body 'http://foo:123@127.0.0.1:3003/x/y/z/a/?json' | jq '.subItems | length')
[ "$cnt" == "0" ] && fail "subItems should not be 0"

# --index-dir-user

cnt=$(curl_get_body 'http://127.0.0.1:3003/x/y/z/b/?json' | jq '.subItems | length')
assert "$cnt" '0'

cnt=$(curl_get_body 'http://baz:789@127.0.0.1:3003/x/y/z/b/?json' | jq '.subItems | length')
assert "$cnt" '0'

cnt=$(curl_get_body 'http://bar:456@127.0.0.1:3003/x/y/z/b/?json' | jq '.subItems | length')
[ "$cnt" == "0" ] && fail "subItems should not be 0"

jobs -p | xargs kill &> /dev/null
