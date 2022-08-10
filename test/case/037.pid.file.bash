#!/bin/bash

cleanup() {
	rm -f "$fs"/vhost1/pid.txt
}

source "$root"/lib.bash

GHFS_PID_FILE="$fs"/vhost1/pid.txt "$ghfs" -l 3003 -r "$fs"/vhost1 &
sleep 0.05 # wait server ready

assert $(cat "$fs"/vhost1/pid.txt) $!

cleanup
jobs -p | xargs kill &> /dev/null
