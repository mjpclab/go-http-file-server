#!/bin/bash

cleanup() {
	rm -f "$fs"/uploaded/[12]/*.tmp
}

source "$root"/lib.bash

"$ghfs" -l 3003 -r "$fs"/uploaded --upload /1 --upload-dir "$fs"/uploaded/2 &
sleep 0.05 # wait server ready
cleanup

content1='uploaded/1/uploaded.tmp'
curl_upload_content 'http://127.0.0.1:3003/1?upload' files "$content1" uploaded.tmp
uploaded1=$(cat "$fs"/uploaded/1/uploaded.tmp)
assert "$uploaded1" "$content1"

content2='uploaded/2/uploaded.tmp'
curl_upload_content 'http://127.0.0.1:3003/2?upload' files "$content2" uploaded.tmp
uploaded2=$(cat "$fs"/uploaded/2/uploaded.tmp)
assert "$uploaded2" "$content2"

cleanup
jobs -p | xargs kill
