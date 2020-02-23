#!/bin/bash

source "$root"/lib.bash

"$ghfs" -l 3003 -r "$fs"/vhost1 -a :world:"$fs"/vhost1/go :/not/exist/name:"$fs"/vhost1/yes &

file1=$(http_get_body 127.0.0.1:3003/file1.txt)
assert "$file1" 'vhost1/file1.txt'

hello=$(http_get_body 127.0.0.1:3003/hello/index.txt)
assert "$hello" 'vhost1/hello/index.txt'

go=$(http_get_body 127.0.0.1:3003/world/index.txt)
assert "$go" 'vhost1/go/index.txt'

yes=$(http_get_body 127.0.0.1:3003/not/exist/name/index.txt)
assert "$yes" 'vhost1/yes/index.txt'

kill %1
