#!/bin/bash

source "$root"/lib.bash

"$ghfs" -l 3003 -r "$fs"/vhost1 -I index.txt &
sleep 0.05 # wait server ready

yes=$(curl_get_body http://127.0.0.1:3003/yes)
assert "$yes" 'vhost1/yes/index.txt'

jobs -p | xargs kill &> /dev/null
