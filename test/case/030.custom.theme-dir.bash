#!/bin/bash

source "$root"/lib.bash

"$ghfs" -l 3003 -r "$fs"/vhost1 --hostname 127.0.0.1 ,,  -l 3003 -r "$fs"/vhost2 --hostname 127.0.0.2 --theme-dir theme-dir &
sleep 0.05 # wait server ready

(curl_get_body 'http://127.0.0.1:3003/' | grep -q -F '<meta http-equiv="comment" content="index.html from custom theme"/>') &&
	fail "Should not use custom theme for vhost 1"

(curl_get_body 'http://127.0.0.1:3003/?asset=index.css' | grep -q -F '/* index.css from custom theme */') &&
	fail "Should not use custom theme for vhost 1"

(curl_get_body 'http://127.0.0.1:3003/?asset=index.js' | grep -q -F '/* index.js from custom theme */') &&
	fail "Should not use custom theme for vhost 1"

(curl_get_body 'http://127.0.0.2:3003/' | grep -q -F '<meta http-equiv="comment" content="index.html from custom theme"/>') ||
	fail "Should use custom theme for vhost 2"

(curl_get_body 'http://127.0.0.2:3003/?asset=index.css' | grep -q -F '/* index.css from custom theme */') ||
	fail "Should use custom theme for vhost 2"

(curl_get_body 'http://127.0.0.2:3003/?asset=index.js' | grep -q -F '/* index.js from custom theme */') ||
	fail "Should use custom theme for vhost 2"

jobs -p | xargs kill
