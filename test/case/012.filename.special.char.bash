#!/bin/bash

source "$root"/lib.bash

"$ghfs" -l 3003 -r "$fs"/vhost1 &
sleep 0.05 # wait server ready

(curl_get_body http://127.0.0.1:3003/special-char/ | grep -q -G '>first<em>\\n</em>second/</') ||
	fail 'first\\nsecond filename not displayed correctly'

(curl_get_body http://127.0.0.1:3003/special-char/ | grep -q -G '>important<em>\\t</em>notice.txt<') ||
	fail 'important\\tnotice.txt filename not displayed correctly'

notice=$(curl_get_body http://127.0.0.1:3003/special-char/important%09notice.txt)
assert "$notice" 'vhost1/special-char/important\tnotice.txt'

kill %1
