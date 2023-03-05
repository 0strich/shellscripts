#!/bin/bash
set -e

# Import Create Direcory functions
source ./create.sh
source ./install.sh

# Default values for options
INSTALL_DOCKER=false
INSTALL_NODE=false
INSTALL_NGINX=false
INSTALL_MONGODB=false

# Usage information
function usage() {
	echo "Usage: $0 [-d] [-n] [-x] [-m]"
	echo "Options:"
	echo "  -d  Install Docker"
	echo "  -n  Install Node.js"
	echo "  -x  Install Nginx"
	echo "  -m  Install MongoDB"
}

# Parse options
while getopts ":dnxm" opt; do
	case ${opt} in
	d)
		INSTALL_DOCKER=true
		;;
	n)
		INSTALL_NODE=true
		;;
	x)
		INSTALL_NODE=true
		;;
	m)
		INSTALL_NODE=true
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

# Check if no options are provided
if ! $INSTALL_DOCKER && ! $INSTALL_NODE && ! $INSTALL_NPM; then
	echo "At least one option is required."
	usage
	exit 1
fi

if $INSTALL_DOCKER; then
	install_docker
fi

if $INSTALL_NODE; then
	install_nodejs
fi

if $INSTALL_NGINX; then
	clone_nginx
	pushd /project/nginx
	./start.sh
fi

if $INSTALL_MONGODB; then
	clone_mongodb
fi

# Final message
echo "Setup complete."
