#!/bin/bash

source "$root"/lib.bash

"$ghfs" -l 3003 -r "$fs"/vhost1 --header foo:bar --header 'set-cookie:name1=value1' --header 'set-cookie:name2=value2' &
sleep 0.05 # wait server ready

(curl_get_header http://127.0.0.1:3003/ | grep -q -i 'foo:\s*bar') ||
	fail "Custom header 'foo:bar' not exists"

(curl_get_header http://127.0.0.1:3003/file1.txt | grep -q -i 'foo:\s*bar') ||
	fail "Custom header 'foo:bar' not exists"

(curl_get_header http://127.0.0.1:3003/ | grep -q -i 'set-cookie:\s*name1=value1') ||
	fail "Custom header 'set-cookie:name1=value1' not exists"

(curl_get_header http://127.0.0.1:3003/ | grep -q -i 'set-cookie:\s*name2=value2') ||
	fail "Custom header 'set-cookie:name2=value2' not exists"

jobs -p | xargs kill
