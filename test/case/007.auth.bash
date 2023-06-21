#!/bin/bash

source "$root"/lib.bash

"$ghfs" -l 3003 -r "$fs"/vhost1 --auth /hello --user alice:AliceSecret -E '' &
sleep 0.05 # wait server ready

yesstatus=$(curl_get_status http://127.0.0.1:3003/yes/)
assert "$yesstatus" '200'

hellostatus=$(curl_get_status http://127.0.0.1:3003/hello/)
assert "$hellostatus" '401'

userhellostatus=$(curl_get_status http://alice:AliceSecret@127.0.0.1:3003/hello/)
assert "$userhellostatus" '200'

userhelloheadstatus=$(curl_head_status http://alice:AliceSecret@127.0.0.1:3003/hello/)
assert "$userhelloheadstatus" '200'

hellostatus=$(curl_get_status http://127.0.0.1:3003/yes/?auth)
assert "$hellostatus" '401'

hellostatus=$(curl_get_status http://alice:AliceSecret@127.0.0.1:3003/yes/?auth)
assert "$hellostatus" '302'

jobs -p | xargs kill &> /dev/null
