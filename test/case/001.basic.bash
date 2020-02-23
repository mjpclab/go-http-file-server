#!/bin/bash

source "$root"/lib.bash

"$ghfs" -l 3003 -r "$fs"/vhost1 &

file1=$(http_get_body 127.0.0.1:3003/file1.txt)
assert "$file1" 'vhost1/file1.txt'

kill %1
