#!/bin/bash

source "$root"/lib.bash

"$ghfs" -l 3003 -r "$fs"/vhost1 --auth /hello --user alice:AliceSecret &

yesstatus=$(http_get_status 127.0.0.1:3003/yes/)
assert "$yesstatus" '200'

hellostatus=$(http_get_status 127.0.0.1:3003/hello/)
assert "$hellostatus" '401'

userhellostatus=$(http_get_status alice:AliceSecret@127.0.0.1:3003/hello/)
assert "$userhellostatus" '200'

userhelloheadstatus=$(http_head_status alice:AliceSecret@127.0.0.1:3003/hello/)
assert "$userhelloheadstatus" '200'

kill %1
