#!/bin/bash

source "$root"/lib.bash

"$ghfs" -l 3003 -r "$fs"/vhost1 -a :shortcut/vhost2:"$fs"/vhost2/ \
	--global-header foo:bar \
	--global-header 'set-cookie:name1=value1' \
	--global-header 'set-cookie:name2=value2' \
	--header '|hello|X-Hello-Name|X-Hello-Value' \
	--header '|/shortcut|X-Sc-Name|X-Sc-Value' \
	--header-dir ":$fs/vhost1/world:X-World-Name:X-World-Value" \
	--header-dir ":$fs/vhost2:X-Vh2-Name:X-Vh2-Value" \
	&
sleep 0.05 # wait server ready

(curl_get_header http://127.0.0.1:3003/ | grep -q -i 'foo:\s*bar') ||
	fail "Custom header 'foo:bar' should exists"

(curl_get_header http://127.0.0.1:3003/file1.txt | grep -q -i 'foo:\s*bar') ||
	fail "Custom header 'foo:bar' should exists"

(curl_get_header http://127.0.0.1:3003/ | grep -q -i 'set-cookie:\s*name2=value2') ||
	fail "Custom header 'set-cookie:name2=value2' should exists"

(curl_get_header http://127.0.0.1:3003/hello | grep -q -i 'X-Hello-Name:\s*X-Hello-Value') ||
	fail "Custom header 'X-Hello-Name:X-Hello-Value' should exists"

(curl_get_header http://127.0.0.1:3003/shortcut/vhost2/file1.txt | grep -q -i 'X-Sc-Name:\s*X-Sc-Value') ||
	fail "Custom header 'X-Sc-Name:X-Sc-Value' should exists"

(curl_get_header http://127.0.0.1:3003/world | grep -q -i 'X-World-Name:\s*X-World-Value') ||
	fail "Custom header 'X-World-Name:X-World-Value' should exists"

(curl_get_header http://127.0.0.1:3003/shortcut/vhost2 | grep -q -i 'X-Vh2-Name:\s*X-Vh2-Value') ||
	fail "Custom header 'X-Vh2-Name:X-Vh2-Value' should exists"

jobs -p | xargs kill &> /dev/null
