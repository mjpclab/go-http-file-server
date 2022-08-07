#!/bin/bash

source "$root"/lib.bash

"$ghfs" -l 3003 -r "$fs"/vhost1 -a :/hello-index.txt:"$fs"/vhost1/hello/index.txt &
sleep 0.05 # wait server ready

line1=$(curl_get_body http://127.0.0.1:3003/ | grep -n './file1.txt' | cut -d ':' -f1)
line2=$(curl_get_body http://127.0.0.1:3003/ | grep -n './file1-1.txt' | cut -d ':' -f1)
line3=$(curl_get_body http://127.0.0.1:3003/ | grep -n './file2.txt' | cut -d ':' -f1)
line4=$(curl_get_body http://127.0.0.1:3003/ | grep -n './file10.txt' | cut -d ':' -f1)
line5=$(curl_get_body http://127.0.0.1:3003/ | grep -n './hello-index.txt' | cut -d ':' -f1)

[ -z "$line1" ] && fail 'line1 is empty'
[ -z "$line2" ] && fail 'line2 is empty'
[ -z "$line3" ] && fail 'line3 is empty'
[ -z "$line4" ] && fail 'line4 is empty'
[ -z "$line5" ] && fail 'line5 is empty'

[ ! "$line1" -lt "$line2" ] && fail "line1 is not before line2"
[ ! "$line2" -lt "$line3" ] && fail "line2 is not before line3"
[ ! "$line3" -lt "$line4" ] && fail "line3 is not before line4"
[ ! "$line4" -lt "$line5" ] && fail "line4 is not before line5"

jobs -p | xargs kill &> /dev/null
