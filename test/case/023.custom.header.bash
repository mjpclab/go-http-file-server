#!/bin/bash

source "$root"/lib.bash

"$ghfs" -l 3003 -r "$fs"/vhost1 --header foo:bar &
sleep 0.05 # wait server ready

(curl_get_header http://127.0.0.1:3003/ | grep -q -i 'foo:\s*bar') ||
	fail "Custom header not exists"

(curl_get_header http://127.0.0.1:3003/file1.txt | grep -q -i 'foo:\s*bar') ||
	fail "Custom header not exists"

kill %1
