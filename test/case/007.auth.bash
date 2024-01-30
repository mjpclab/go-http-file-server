#!/bin/bash

source "$root"/lib.bash

"$ghfs" -l 3003 -r "$fs"/vhost1 --auth /hello --user alice:AliceSecret -E '' &
sleep 0.05 # wait server ready

status=$(curl_get_status http://127.0.0.1:3003/yes/)
assert "$status" '200'

status=$(curl_get_status http://127.0.0.1:3003/hello/)
assert "$status" '401'

status=$(curl_get_status http://alice:AliceSecret@127.0.0.1:3003/hello/)
assert "$status" '200'

status=$(curl_head_status http://alice:AliceSecret@127.0.0.1:3003/hello/)
assert "$status" '200'

status=$(curl_get_status http://127.0.0.1:3003/yes/?auth)
assert "$status" '401'

status=$(curl_get_status http://bob:BobSecret@127.0.0.1:3003/yes/)
assert "$status" '200'

status=$(curl_get_status http://alice:AliceSecret@127.0.0.1:3003/yes/?auth)
assert "$status" '302'

jobs -p | xargs kill &> /dev/null
