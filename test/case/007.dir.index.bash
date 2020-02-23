#!/bin/bash

source "$root"/lib.bash

"$ghfs" -l 3003 -r "$fs"/vhost1 -I index.txt &

yes=$(http_get_body 127.0.0.1:3003/yes)
assert "$yes" 'vhost1/yes/index.txt'

kill %1
