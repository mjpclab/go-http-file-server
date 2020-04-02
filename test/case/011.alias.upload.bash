#!/bin/bash

cleanup() {
	rm -f "$fs"/uploaded/2/*.tmp
}

source "$root"/lib.bash

"$ghfs" -l 3003 -r "$fs"/vhost1 --alias :/my/upload:"$fs"/uploaded/2 --upload /my/upload &
sleep 0.05 # wait server ready
cleanup

content='my/upload/uploaded.tmp'
curl_upload_content 'http://127.0.0.1:3003/my/upload?upload' files "$content" uploaded.tmp
uploaded=$(cat "$fs"/uploaded/2/uploaded.tmp)
assert "$uploaded" "$content"

cleanup
kill %1
