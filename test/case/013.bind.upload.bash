#!/bin/bash

cleanup() {
	rm -f "$fs"/uploaded/2/*.tmp
	rm -fr "$fs"/vhost1/[Mm][Yy]*
}

source "$root"/lib.bash

"$ghfs" -l 3003 -r "$fs"/vhost1 --bind :/my/upload:"$fs"/uploaded/2 --upload / --mkdir / &> /dev/null &
sleep 0.05 # wait server ready
cleanup

content='my/upload/uploaded.tmp'
curl_upload_content 'http://127.0.0.1:3003/my/upload?upload' file "$content" uploaded.tmp
uploaded=$(cat "$fs"/uploaded/2/uploaded.tmp)
assert "$uploaded" "$content"

content='MY/Upload/uploaded2.tmp'
curl_upload_content 'http://127.0.0.1:3003/my/upload?upload' file "$content" 'temp/dir/uploaded2.tmp'
uploaded=$(cat "$fs"/uploaded/2/uploaded2.tmp)
assert "$uploaded" "$content"

curl_upload_content 'http://127.0.0.1:3003/?upload' file mycontent My
[ -e "$fs"/vhost1/My ] && fail "$fs/vhost1/My should not exists"

curl_upload_content 'http://127.0.0.1:3003/?upload' file mycontent 'dir/to/My'
[ -e "$fs"/vhost1/My ] && fail "$fs/vhost1/My should not exists"

curl_upload_content 'http://127.0.0.1:3003/?upload' dirfile mycontent 'My/Myfile'
[ -e "$fs"/vhost1/My ] && fail "$fs/vhost1/My should not exists"

curl_upload_content 'http://127.0.0.1:3003/?upload' dirfile mycontent 'My/Mydir/file'
[ -e "$fs"/vhost1/My ] && fail "$fs/vhost1/My should not exists"

cleanup
jobs -p | xargs kill
