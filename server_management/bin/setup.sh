#!/bin/bash
set -e

# Import Create Direcory functions
source ./create_directory.sh

# Default values for options
INSTALL_DOCKER=false
INSTALL_NODE=false

# Usage information
function usage() {
	echo "Usage: $0 [-d] [-n]"
	echo "Options:"
	echo "  -d  Install Docker"
	echo "  -n  Install Node.js"
}

# Parse options
while getopts ":dn" opt; do
	case ${opt} in
	d)
		INSTALL_DOCKER=true
		;;
	n)
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

# Functions for installation
function install_docker() {
	echo "Installing Docker..."
	# apt 업데이트, docker 설치
	apt update -y && apt install docker.io -y

	# JSOB processer 설치
	apt-get install jq -y

	# VERSION, DESTINATINO 환경변수 설정
	LATEST_VERSION=$(curl --silent https://api.github.com/repos/docker/compose/releases/1.29.2 | jq .name -r)
	VERSION=$(curl --silent https://api.github.com/repos/docker/compose/releases/1.29.2 | jq .name -r)
	DESTINATION=/usr/local/bin/docker-compose

	# docker-compose 설치
	sudo curl -L https://github.com/docker/compose/releases/download/1.29.2/docker-compose-$(uname -s)-$(uname -m) -o $DESTINATION

	# DESTINATION 경로 권한 변경
	sudo chmod 755 $DESTINATION
	echo "Docker installed."
}

function install_nodejs() {
	echo "Installing Node.js..."
	curl -sL https://deb.nodesource.com/setup_14.x | sudo -E bash -
	sudo apt-get install -y nodejs

	echo "Installing NVM"
	sudo apt install -y curl build-essential libssl-dev
	curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.1/install.sh | bash
	soruce ~/.bashrc

	echo "Node.js installed."
}

# Install packages based on options
if $INSTALL_DOCKER; then
	install_docker
fi

if $INSTALL_NODE; then
	install_nodejs
fi

# Final message
echo "Setup complete."
