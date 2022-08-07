#!/bin/bash

source "$root"/lib.bash

"$ghfs" -l 3003 -r "$fs"/vhost1 --to-https &
sleep 0.05 # wait server ready

status=$(curl_get_status http://127.0.0.1:3003/)
assert "$status" '200'

jobs -p | xargs kill &> /dev/null


"$ghfs" --to-https 3005 --listen-plain 3003 --listen-tls 3004 --listen-tls 3005 -c "$cert"/localhost.crt -k "$cert"/localhost.key -r "$fs"/vhost1 &
sleep 0.05 # wait server ready

status=$(curl_get_status http://127.0.0.1:3003/)
assert "$status" '301'

(curl_get_header http://127.0.0.1:3003/ | grep -i -q 'Location:\s*https://127.0.0.1:3005/') || fail 'Location not found'

jobs -p | xargs kill &> /dev/null


# --to-https empty
"$ghfs" --to-https --listen-plain 3003 --listen-tls 3004 -c "$cert"/localhost.crt -k "$cert"/localhost.key -r "$fs"/vhost1 &
sleep 0.05 # wait server ready

status=$(curl_get_status http://127.0.0.1:3003/)
assert "$status" '301'

jobs -p | xargs kill &> /dev/null


"$ghfs" --to-https --listen-plain 3003 --listen-tls :3004 -c "$cert"/localhost.crt -k "$cert"/localhost.key -r "$fs"/vhost1 &
sleep 0.05 # wait server ready

status=$(curl_get_status http://127.0.0.1:3003/)
assert "$status" '301'

jobs -p | xargs kill &> /dev/null


"$ghfs" --to-https --listen-plain 3003 --listen-tls 127.0.0.1:3004 -c "$cert"/localhost.crt -k "$cert"/localhost.key -r "$fs"/vhost1 &
sleep 0.05 # wait server ready

status=$(curl_get_status http://127.0.0.1:3003/)
assert "$status" '301'

jobs -p | xargs kill &> /dev/null

# --to-https port
"$ghfs" --to-https 3004 --listen-plain 3003 --listen-tls 3004 -c "$cert"/localhost.crt -k "$cert"/localhost.key -r "$fs"/vhost1 &
sleep 0.05 # wait server ready

status=$(curl_get_status http://127.0.0.1:3003/)
assert "$status" '301'

jobs -p | xargs kill &> /dev/null


"$ghfs" --to-https 3004 --listen-plain 3003 --listen-tls :3004 -c "$cert"/localhost.crt -k "$cert"/localhost.key -r "$fs"/vhost1 &
sleep 0.05 # wait server ready

status=$(curl_get_status http://127.0.0.1:3003/)
assert "$status" '301'

jobs -p | xargs kill &> /dev/null


"$ghfs" --to-https 3004 --listen-plain 3003 --listen-tls 127.0.0.1:3004 -c "$cert"/localhost.crt -k "$cert"/localhost.key -r "$fs"/vhost1 &
sleep 0.05 # wait server ready

status=$(curl_get_status http://127.0.0.1:3003/)
assert "$status" '301'

jobs -p | xargs kill &> /dev/null

# --to-https :port
"$ghfs" --to-https :3004 --listen-plain 3003 --listen-tls 3004 -c "$cert"/localhost.crt -k "$cert"/localhost.key -r "$fs"/vhost1 &
sleep 0.05 # wait server ready

status=$(curl_get_status http://127.0.0.1:3003/)
assert "$status" '301'

jobs -p | xargs kill &> /dev/null


"$ghfs" --to-https :3004 --listen-plain 3003 --listen-tls :3004 -c "$cert"/localhost.crt -k "$cert"/localhost.key -r "$fs"/vhost1 &
sleep 0.05 # wait server ready

status=$(curl_get_status http://127.0.0.1:3003/)
assert "$status" '301'

jobs -p | xargs kill &> /dev/null


"$ghfs" --to-https :3004 --listen-plain 3003 --listen-tls 127.0.0.1:3004 -c "$cert"/localhost.crt -k "$cert"/localhost.key -r "$fs"/vhost1 &
sleep 0.05 # wait server ready

status=$(curl_get_status http://127.0.0.1:3003/)
assert "$status" '301'

jobs -p | xargs kill &> /dev/null

# --to-https IPv4:port
"$ghfs" --to-https 127.0.0.1:3004 --listen-plain 3003 --listen-tls 3004 -c "$cert"/localhost.crt -k "$cert"/localhost.key -r "$fs"/vhost1 &
sleep 0.05 # wait server ready

status=$(curl_get_status http://127.0.0.1:3003/)
assert "$status" '301'

jobs -p | xargs kill &> /dev/null


"$ghfs" --to-https 127.0.0.1:3004 --listen-plain 3003 --listen-tls :3004 -c "$cert"/localhost.crt -k "$cert"/localhost.key -r "$fs"/vhost1 &
sleep 0.05 # wait server ready

status=$(curl_get_status http://127.0.0.1:3003/)
assert "$status" '301'

jobs -p | xargs kill &> /dev/null


"$ghfs" --to-https 127.0.0.1:3004 --listen-plain 3003 --listen-tls 127.0.0.1:3004 -c "$cert"/localhost.crt -k "$cert"/localhost.key -r "$fs"/vhost1 &
sleep 0.05 # wait server ready

status=$(curl_get_status http://127.0.0.1:3003/)
assert "$status" '301'

jobs -p | xargs kill &> /dev/null

# --to-https IPv6:port
"$ghfs" --to-https '[::1]:3004' --listen-plain 3003 --listen-tls 3004 -c "$cert"/localhost.crt -k "$cert"/localhost.key -r "$fs"/vhost1 &
sleep 0.05 # wait server ready

status=$(curl_get_status http://127.0.0.1:3003/)
assert "$status" '301'

jobs -p | xargs kill &> /dev/null


"$ghfs" --to-https '[::1]:3004' --listen-plain 3003 --listen-tls :3004 -c "$cert"/localhost.crt -k "$cert"/localhost.key -r "$fs"/vhost1 &
sleep 0.05 # wait server ready

status=$(curl_get_status http://127.0.0.1:3003/)
assert "$status" '301'

jobs -p | xargs kill &> /dev/null


"$ghfs" --to-https '[::1]:3004' --listen-plain 3003 --listen-tls 127.0.0.1:3004 -c "$cert"/localhost.crt -k "$cert"/localhost.key -r "$fs"/vhost1 &
sleep 0.05 # wait server ready

status=$(curl_get_status http://127.0.0.1:3003/)
assert "$status" '301'

jobs -p | xargs kill &> /dev/null

# --to-https bad port
"$ghfs" --to-https 3005 --listen-plain 3003 --listen-tls 3004 -c "$cert"/localhost.crt -k "$cert"/localhost.key -r "$fs"/vhost1 &
sleep 0.05 # wait server ready

status=$(curl_get_status http://127.0.0.1:3003/)
assert "$status" '200'

jobs -p | xargs kill &> /dev/null


"$ghfs" --to-https 3005 --listen-plain 3003 --listen-tls :3004 -c "$cert"/localhost.crt -k "$cert"/localhost.key -r "$fs"/vhost1 &
sleep 0.05 # wait server ready

status=$(curl_get_status http://127.0.0.1:3003/)
assert "$status" '200'

jobs -p | xargs kill &> /dev/null


"$ghfs" --to-https 3005 --listen-plain 3003 --listen-tls 127.0.0.1:3004 -c "$cert"/localhost.crt -k "$cert"/localhost.key -r "$fs"/vhost1 &
sleep 0.05 # wait server ready

status=$(curl_get_status http://127.0.0.1:3003/)
assert "$status" '200'

jobs -p | xargs kill &> /dev/null
