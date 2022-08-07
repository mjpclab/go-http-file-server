#!/bin/bash

source "$root"/lib.bash

"$ghfs" -l 3003 -R -a :world:"$fs"/vhost1/go :/not/exist/name:"$fs"/vhost1/yes &
sleep 0.05 # wait server ready

file1status=$(curl_get_status http://127.0.0.1:3003/file1.txt)
assert "$file1status" '404'

file1headstatus=$(curl_head_status http://127.0.0.1:3003/file1.txt)
assert "$file1headstatus" '404'

hellostatus=$(curl_get_status http://127.0.0.1:3003/hello/index.txt)
assert "$hellostatus" '404'

helloheadstatus=$(curl_head_status http://127.0.0.1:3003/hello/index.txt)
assert "$helloheadstatus" '404'

go=$(curl_get_body http://127.0.0.1:3003/world/index.txt)
assert "$go" 'vhost1/go/index.txt'

yes=$(curl_get_body http://127.0.0.1:3003/not/exist/name/index.txt)
assert "$yes" 'vhost1/yes/index.txt'

jobs -p | xargs kill &> /dev/null
