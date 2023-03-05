#!/bin/bash
set -e

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
	source ~/.bashrc

	echo "Node.js installed."
}

function install_nginx() {
	echo "Installing Nginx..."
	sudo apt-get update -y
	sudo apt-get install -y nginx
	echo "Nginx installed."
}

function install_mongodb() {
	echo "Installing MongoDB..."
	sudo apt-get update -y
	sudo apt-get install -y mongodb
	echo "MongoDB installed."
}
