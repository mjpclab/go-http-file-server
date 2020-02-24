#!/bin/bash

source "$root"/lib.bash

"$ghfs" -l 3003 -r "$fs"/vhost1 -I index.txt -a :/hello/index.txt:"$fs"/vhost1/world/index.txt &

world=$(http_get_body 127.0.0.1:3003/hello)
assert "$world" 'vhost1/world/index.txt'

kill %1
