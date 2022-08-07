#!/bin/bash

cleanup() {
	rm -rf "$fs"/uploaded/[12]/{*.tmp,*/}
}

source "$root"/lib.bash

"$ghfs" -l 3003 -r "$fs"/uploaded --upload /1 --upload-dir "$fs"/uploaded/2 --mkdir /2 -E '' &
sleep 0.05 # wait server ready
cleanup

# file upload to /1 - valid
content='uploaded/1/uploaded.tmp'
curl_upload_content 'http://127.0.0.1:3003/1?upload' file "$content" uploaded.tmp
uploaded=$(cat "$fs"/uploaded/1/uploaded.tmp)
assert "$uploaded" "$content"

# file upload to /2 - valid
content='uploaded/2/uploaded.tmp'
curl_upload_content 'http://127.0.0.1:3003/2?upload' file "$content" uploaded.tmp
uploaded=$(cat "$fs"/uploaded/2/uploaded.tmp)
assert "$uploaded" "$content"


# dir file upload to /1 - INVALID
content='upload/1/dir/file'
curl_upload_content 'http://127.0.0.1:3003/1?upload' dirfile "$content" dir/uploaded.tmp
[ ! -e "$fs"/uploaded/1/dir/uploaded.tmp ] || fail "/uploaded/1/dir/uploaded.tmp should not be exists"

# dir file upload to /2 - valid
content='upload/2/dir/file'
curl_upload_content 'http://127.0.0.1:3003/2?upload' dirfile "$content" dir/uploaded.tmp
[ -e "$fs"/uploaded/2/dir/uploaded.tmp ] || fail "/uploaded/2/dir/uploaded.tmp should be exists"
uploaded=$(cat "$fs"/uploaded/2/dir/uploaded.tmp)
assert "$uploaded" "$content"

# sub dir file upload to /2 - valid
content='upload/2/dir/sub/file'
curl_upload_content 'http://127.0.0.1:3003/2?upload' dirfile "$content" dir/sub/uploaded.tmp
[ -e "$fs"/uploaded/2/dir/sub/uploaded.tmp ] || fail "/uploaded/2/dir/sub/uploaded.tmp should be exists"
uploaded=$(cat "$fs"/uploaded/2/dir/sub/uploaded.tmp)
assert "$uploaded" "$content"


# inner dir file upload to /1 - valid
content='upload/1/idir/file'
curl_upload_content 'http://127.0.0.1:3003/1?upload' innerdirfile "$content" idir/iuploaded.tmp
uploaded=$(cat "$fs"/uploaded/1/iuploaded.tmp)
assert "$uploaded" "$content"

# inner sub dir file upload to /1 - INVALID
content='upload/1/idir/isub/file'
curl_upload_content 'http://127.0.0.1:3003/1?upload' innerdirfile "$content" idir/isub/iuploaded.tmp
[ ! -e "$fs"/uploaded/2/isub/iuploaded.tmp ] || fail "/uploaded/2/isub/iuploaded.tmp should not be exists"

# inner dir file upload to /2 - valid
content='upload/2/idir/file'
curl_upload_content 'http://127.0.0.1:3003/2?upload' innerdirfile "$content" idir/iuploaded.tmp
uploaded=$(cat "$fs"/uploaded/2/iuploaded.tmp)
assert "$uploaded" "$content"

# inner sub dir file upload to /2 - valid
content='upload/2/idir/isub/file'
curl_upload_content 'http://127.0.0.1:3003/2?upload' innerdirfile "$content" idir/isub/iuploaded.tmp
[ -e "$fs"/uploaded/2/isub/iuploaded.tmp ] || fail "/uploaded/2/isub/iuploaded.tmp should be exists"
uploaded=$(cat "$fs"/uploaded/2/isub/iuploaded.tmp)
assert "$uploaded" "$content"

cleanup
jobs -p | xargs kill &> /dev/null
