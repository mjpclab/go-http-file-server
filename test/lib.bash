#!/bin/bash

assert() {
	expect="$2"
	actual="$1"
	if [ "$expect" != "$actual" ]; then
		echo -e "$(basename $0): expect \"\e[0;32m$expect\e[0m\", got \"\e[1;31m$actual\e[0m\""
	fi;
}

http_head_status() {
	url=http://"$1"
	curl -s -I "$url" | head -n 1 | cut -d ' ' -f 2
}

http_get_status() {
	url=http://"$1"
	curl -s -i "$url" | head -n 1 | cut -d ' ' -f 2
}

http_get_body() {
	url=http://"$1"
	curl -s "$url"
}

https_head_status() {
	url=https://"$1"
	curl -s -k -I "$url" | head -n 1 | cut -d ' ' -f 2
}

https_get_status() {
	url=https://"$1"
	curl -s -k -i "$url" | head -n 1 | cut -d ' ' -f 2
}

https_get_body() {
	url=https://"$1"
	curl -s -k "$url"
}
