type=$1
message=$2

push () {
	npm run $type
	git add .
	git commit -m "$message"
	git push
}

case $1 in
	patch) push;;
	minor) push;;
	major) push;;
esac

