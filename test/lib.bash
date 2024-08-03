#!/bin/bash

assert() {
	local expect="$2"
	local actual="$1"
	if [ "$expect" != "$actual" ]; then
		echo -e "$(basename $0):${BASH_LINENO[0]} expect \"\e[0;32m$expect\e[0m\", got \"\e[1;31m$actual\e[0m\"" >&2
	fi
}

fail() {
	local msg="$1"
	echo -e "$(basename $0):${BASH_LINENO[0]} \e[1;31m$msg\e[0m" >&2
}

curl_head_status() {
	local args=("$@")
	local urlindex=$[ ${#args[@]} - 1 ]
	local opts=("${args[@]:0:urlindex}")
	local url="${args[urlindex]}"

	curl -s -k -I "${opts[@]}" "$url" | head -n 1 | cut -d ' ' -f 2
}

curl_get_status() {
	local args=("$@")
	local urlindex=$[ ${#args[@]} - 1 ]
	local opts=("${args[@]:0:urlindex}")
	local url="${args[urlindex]}"

	curl -s -k -i "${opts[@]}" "$url" | head -n 1 | cut -d ' ' -f 2
}

curl_get_header() {
	local args=("$@")
	local urlindex=$[ ${#args[@]} - 1 ]
	local opts=("${args[@]:0:urlindex}")
	local url="${args[urlindex]}"

	curl -s -k -i "${opts[@]}" "$url" | sed -e '/^$/q'
}

curl_get_body() {
	local args=("$@")
	local urlindex=$[ ${#args[@]} - 1 ]
	local opts=("${args[@]:0:urlindex}")
	local url="${args[urlindex]}"

	curl -s -k "${opts[@]}" "$url"
}

curl_post_status() {
	curl_get_status -X POST "$@"
}

curl_upload_content() {
	local url="$1"
	local name="$2"
	local value="$3"
	local filename="$4"
	curl -s -k -F "$name=$value;filename=$filename" "$url"
}
