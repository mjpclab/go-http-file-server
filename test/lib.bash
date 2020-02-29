#!/bin/bash

assert() {
	expect="$2"
	actual="$1"
	if [ "$expect" != "$actual" ]; then
		echo -e "$(basename $0): expect \"\e[0;32m$expect\e[0m\", got \"\e[1;31m$actual\e[0m\"" >&2
	fi
}

fail() {
	msg="$1"
	echo -e "$(basename $0): \e[1;31m$msg\e[0m" >&2
}

curl_head_status() {
	url="$1"
	curl -s -k -I "$url" | head -n 1 | cut -d ' ' -f 2
}

curl_get_status() {
	url="$1"
	opts="$2"
	curl -s -k -i "$url" | head -n 1 | cut -d ' ' -f 2
}

curl_get_body() {
	url="$1"
	opts="$2"
	curl -s -k "$url"
}

curl_upload_content() {
	url="$1"
	name="$2"
	value="$3"
	filename="$4"
	curl -s -k -F "$name=$value;filename=$filename" "$url"
}
