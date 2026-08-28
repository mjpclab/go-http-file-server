#!/bin/bash

source "$root"/lib.bash

"$ghfs" -l 3003 -r "$fs"/vhost1 -a :/x/y/z:"$fs"/vhost2 --user foo:123 bar:456 --global-list false --list /hello --list-dir "$fs"/vhost1/world --list-user :/x/y/z/a:foo --list-dir-user :"$fs"/vhost2/b:bar -E '' &
sleep 0.05 # wait server ready

# --list

cnt=$(curl_get_body -H 'accept: application/json' 'http://127.0.0.1:3003/' | jq '.subItems | length')
assert "$cnt" '0'

cnt=$(curl_get_body -H 'accept: application/json' 'http://127.0.0.1:3003/go/' | jq '.subItems | length')
assert "$cnt" '0'

cnt=$(curl_get_body -H 'accept: application/json' 'http://127.0.0.1:3003/hello/' | jq '.subItems | length')
[ "$cnt" == "0" ] && fail "subItems should not be 0"

# --list-dir

cnt=$(curl_get_body -H 'accept: application/json' 'http://127.0.0.1:3003/world/' | jq '.subItems | length')
[ "$cnt" == "0" ] && fail "subItems should not be 0"

# --list-user

cnt=$(curl_get_body -H 'accept: application/json' 'http://127.0.0.1:3003/x/y/z/a/' | jq '.subItems | length')
assert "$cnt" '0'

cnt=$(curl_get_body -H 'accept: application/json' 'http://baz:789@127.0.0.1:3003/x/y/z/a/' | jq '.subItems | length')
assert "$cnt" '0'

cnt=$(curl_get_body -H 'accept: application/json' 'http://foo:123@127.0.0.1:3003/x/y/z/a/' | jq '.subItems | length')
[ "$cnt" == "0" ] && fail "subItems should not be 0"

# --list-dir-user

cnt=$(curl_get_body -H 'accept: application/json' 'http://127.0.0.1:3003/x/y/z/b/' | jq '.subItems | length')
assert "$cnt" '0'

cnt=$(curl_get_body -H 'accept: application/json' 'http://baz:789@127.0.0.1:3003/x/y/z/b/' | jq '.subItems | length')
assert "$cnt" '0'

cnt=$(curl_get_body -H 'accept: application/json' 'http://bar:456@127.0.0.1:3003/x/y/z/b/' | jq '.subItems | length')
[ "$cnt" == "0" ] && fail "subItems should not be 0"

jobs -p | xargs kill &> /dev/null
