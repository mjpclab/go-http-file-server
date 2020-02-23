#!/bin/bash

source "$root"/lib.bash

"$ghfs" -l 3003 -r "$fs"/vhost1 --hostname=127.0.0.1 ,,  -l 3003 -r "$fs"/vhost2 --hostname=127.0.0.2 &

vh1file1=$(http_get_body 127.0.0.1:3003/file1.txt)
assert "$vh1file1" 'vhost1/file1.txt'

vh2file1=$(http_get_body 127.0.0.2:3003/file1.txt)
assert "$vh2file1" 'vhost2/file1.txt'

kill %1
