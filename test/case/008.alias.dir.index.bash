#!/bin/bash

source "$root"/lib.bash

"$ghfs" -l 3003 -r "$fs"/vhost1 -I index.txt -a :/foo/bar:"$fs"/vhost2/ -a :/foo/index.txt:"$fs"/vhost1/go/index.txt -E '' &
sleep 0.05 # wait server ready

world=$(curl_get_body http://127.0.0.1:3003/foo)
assert "$world" 'vhost1/go/index.txt'

kill %1
