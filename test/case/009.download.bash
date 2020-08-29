#!/bin/bash

source "$root"/lib.bash

"$ghfs" -l 3003 -r "$fs"/vhost1 &
sleep 0.05 # wait server ready

(curl_get_header http://127.0.0.1:3003/file1.txt | grep -q -i -F 'content-disposition') &&
	fail "Content-Disposition header should not exists"

(curl_get_header http://127.0.0.1:3003/file1.txt?download | grep -q -i -F 'content-disposition') ||
	fail "Content-Disposition header is not exists"

kill %1
