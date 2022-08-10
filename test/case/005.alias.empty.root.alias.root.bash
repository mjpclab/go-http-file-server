#!/bin/bash

source "$root"/lib.bash

"$ghfs" -l 3003 -R -a :/:"$fs"/vhost1 &
sleep 0.05 # wait server ready

file1status=$(curl_get_status http://127.0.0.1:3003/file1.txt)
assert "$file1status" '200'

file1headstatus=$(curl_head_status http://127.0.0.1:3003/file1.txt)
assert "$file1headstatus" '200'

jobs -p | xargs kill &> /dev/null
