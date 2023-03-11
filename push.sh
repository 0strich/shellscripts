#!/bin/bash
set -e

# Default values for options
PUSH_NODE=false

type=$1
message=$2

push() {
	npm run $1
	git add .
	git commit -m "$2"
	git push
}

# Usage information
function usage() {
	echo "Usage: $0 [-n] patch|minor|major commitmessage"
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
if $PUSH_NODE || [[ $1 == patch || $1 == minor || $1 == major ]]; then
	case $1 in
	patch | minor | major)
		echo 1 $1
		echo 2 $2
		# push $1 "$2"
		;;
	*)
		echo "Invalid option: $1" 1>&2
		usage
		exit 1
		;;
	esac
else
	# main
	git add .
	git commit -m "$1"
	git push
fi

# Final message
echo "Complete."
