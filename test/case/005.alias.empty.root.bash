#!/bin/bash

source "$root"/lib.bash

"$ghfs" -l 3003 -R -a :world:"$fs"/vhost1/go :/not/exist/name:"$fs"/vhost1/yes &

file1status=$(http_get_status  127.0.0.1:3003/file1.txt | head -n 1 | cut -d ' ' -f 2)
assert "$file1status" '404'

file1headstatus=$(http_head_status  127.0.0.1:3003/file1.txt | head -n 1 | cut -d ' ' -f 2)
assert "$file1headstatus" '404'

hellostatus=$(http_get_status 127.0.0.1:3003/hello/index.txt | head -n 1 | cut -d ' ' -f 2)
assert "$hellostatus" '404'

helloheadstatus=$(http_head_status 127.0.0.1:3003/hello/index.txt | head -n 1 | cut -d ' ' -f 2)
assert "$helloheadstatus" '404'

go=$(http_get_body 127.0.0.1:3003/world/index.txt)
assert "$go" 'vhost1/go/index.txt'

yes=$(http_get_body 127.0.0.1:3003/not/exist/name/index.txt)
assert "$yes" 'vhost1/yes/index.txt'

kill %1
