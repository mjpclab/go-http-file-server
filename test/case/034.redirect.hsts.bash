#!/bin/bash

source "$root"/lib.bash

"$ghfs" -l 3003 -r "$fs"/vhost1 --hsts &
sleep 0.05 # wait server ready

status=$(curl_get_status http://127.0.0.1:3003/)
assert "$status" '200'

jobs -p | xargs kill


"$ghfs" --hsts --listen-plain 3003 --listen-tls 3004 -c "$cert"/localhost.crt -k "$cert"/localhost.key -r "$fs"/vhost1 &
sleep 0.05 # wait server ready

status=$(curl_get_status http://127.0.0.1:3003/)
assert "$status" '200'

jobs -p | xargs kill
