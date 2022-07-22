#!/bin/bash

assert() {
	expect="$2"
	actual="$1"
	if [ "$expect" != "$actual" ]; then
		echo -e "$(basename $0):${BASH_LINENO[0]} expect \"\e[0;32m$expect\e[0m\", got \"\e[1;31m$actual\e[0m\"" >&2
	fi
}

fail() {
	msg="$1"
	echo -e "$(basename $0):${BASH_LINENO[0]} \e[1;31m$msg\e[0m" >&2
}

curl_head_status() {
	args=($@)
	urlindex=$[ ${#args[@]} - 1 ]
	opts="${args[@]:0:urlindex}"
	url="${args[urlindex]}"

	curl -s -k -I $opts "$url" | head -n 1 | cut -d ' ' -f 2
}

curl_get_status() {
	args=($@)
	urlindex=$[ ${#args[@]} - 1 ]
	opts="${args[@]:0:urlindex}"
	url="${args[urlindex]}"

	curl -s -k -i $opts "$url" | head -n 1 | cut -d ' ' -f 2
}

curl_get_header() {
	args=($@)
	urlindex=$[ ${#args[@]} - 1 ]
	opts="${args[@]:0:urlindex}"
	url="${args[urlindex]}"

	curl -s -k -i $opts "$url" | sed -e '/^$/q'
}

curl_get_body() {
	args=($@)
	urlindex=$[ ${#args[@]} - 1 ]
	opts="${args[@]:0:urlindex}"
	url="${args[urlindex]}"

	curl -s -k $opts "$url"
}

curl_upload_content() {
	url="$1"
	name="$2"
	value="$3"
	filename="$4"
	curl -s -k -F "$name=$value;filename=$filename" "$url"
}
