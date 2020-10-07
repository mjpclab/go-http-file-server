#!/bin/bash

cleanup() {
	rm -f "$fs"/downloaded/*.tmp
}

source "$root"/lib.bash

"$ghfs" -l 3003 -r "$fs"/vhost1 &
sleep 0.05 # wait server ready
cleanup

curl_get_body 'http://127.0.0.1:3003/?sort=/N' | grep -q -F 'go/?sort=/N"'
if [ $? -ne 0 ]; then
	fail "unexpected url param for sort"
fi

cleanup
jobs -p | xargs kill
