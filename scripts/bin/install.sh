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

# Install Go Language
function install_golang() {
	echo "Installing Go Language..."
	wget https://golang.org/dl/go1.18.linux-amd64.tar.gz
	sudo tar -C /usr/local -xzf go1.18.linux-amd64.tar.gz
	rm -rf go1.18.linux-amd64.tar.gz
	echo "export PATH=$PATH:/usr/local/go/bin" >>~/.bashrc
	echo "source ~/.bashrc"
	echo "Go Language installed successfully"
}

function install_nodejs() {
	echo "Installing Node.js..."
	curl -sL https://deb.nodesource.com/setup_16.x | sudo -E bash -
	sudo apt-get install -y nodejs

	echo "Installing NVM"
	sudo apt install -y curl build-essential libssl-dev
	curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.1/install.sh | bash
	echo "source ~/.bashrc"

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

function install_ipfs() {
	# apt 업데이트 & vim 설치
	apt update -y && apt install vim net-tools -y

	mkdir /project && cd /project

	# go설치 & 압축 해제 & 경로 설정 & 적용
	wget https://golang.org/dl/go1.18.2.linux-amd64.tar.gz
	sudo tar -C /usr/local -xzf go1.18.2.linux-amd64.tar.gz
	rm go1.18.2.linux-amd64.tar.gz
	echo "export PATH=$PATH:/usr/local/go/bin" >>~/.bashrc
	echo "source ~/.bashrc"

	# go-ipfs 설치 & ipfs 명령어 폴더 이동 & 경로 설정
	wget https://dist.ipfs.io/go-ipfs/v0.6.0/go-ipfs_v0.6.0_linux-amd64.tar.gz
	rm go-ipfs_v0.6.0_linux-amd64.tar.gz
	sudo mv go-ipfs/ipfs /usr/bin/ipfs
	rm -rf ./go-ipfs
	echo "IPFS_PATH=~/.ipfs" >>~/.bashrc
	echo "source ~/.bashrc"

	# ipfs 초기화
	ipfs init
	ipfs daemon &
}
