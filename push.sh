#!/bin/bash
set -e

# Default values for options
PUSH_NODE=false

type=$1
message=$2

push() {
	npm run $type
	git add .
	git commit -m "$message"
	git push
}

# Usage information
function usage() {
	echo "Usage: $0 [-n]"
	echo "Options:"
	echo "  -n  Push Node.js"
}

# Parse options
while getopts ":n" opt; do
	case ${opt} in
	n)
		PUSH_NODE=true
		;;
	\?)
		usage
		exit 1
		;;
	:)
		echo "Invalid option: -$OPTARG requires an argument" 1>&2
		usage
		exit 1
		;;
	esac
done
shift $((OPTIND - 1))

# Push Node.js
if $PUSH_NODE; then
	case $1 in
	patch | minor | major)
		push
		;;
	*)
		echo "Invalid option: $1" 1>&2
		usage
		exit 1
		;;
	esac
fi

# main
git add .
git commit -m "$1"
git pusk

# Final message
echo "Complete."
