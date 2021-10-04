#!/bin/bash

cleanup() {
	rm -f "$fs"/uploaded/2/*.tmp
	rm -fr "$fs"/vhost1/my*
}

source "$root"/lib.bash

"$ghfs" -l 3003 -r "$fs"/vhost1 --alias :/my/upload:"$fs"/uploaded/2 --upload / --mkdir / &> /dev/null &
sleep 0.05 # wait server ready
cleanup

content='my/upload/uploaded.tmp'
curl_upload_content 'http://127.0.0.1:3003/my/upload?upload' file "$content" uploaded.tmp
uploaded=$(cat "$fs"/uploaded/2/uploaded.tmp)
assert "$uploaded" "$content"

content='my/upload/uploaded2.tmp'
curl_upload_content 'http://127.0.0.1:3003/my/upload?upload' file "$content" 'temp/dir/uploaded2.tmp'
uploaded=$(cat "$fs"/uploaded/2/uploaded2.tmp)
assert "$uploaded" "$content"

curl_upload_content 'http://127.0.0.1:3003/?upload' file mycontent my
[ -e "$fs"/vhost1/my ] && fail "$fs/vhost1/my should not exists"

curl_upload_content 'http://127.0.0.1:3003/?upload' file mycontent 'dir/to/my'
[ -e "$fs"/vhost1/my ] && fail "$fs/vhost1/my should not exists"

curl_upload_content 'http://127.0.0.1:3003/?upload' dirfile mycontent 'my/myfile'
[ -e "$fs"/vhost1/my ] && fail "$fs/vhost1/my should not exists"

curl_upload_content 'http://127.0.0.1:3003/?upload' dirfile mycontent 'my/mydir/file'
[ -e "$fs"/vhost1/my ] && fail "$fs/vhost1/my should not exists"

cleanup
jobs -p | xargs kill
