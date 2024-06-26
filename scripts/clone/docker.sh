#!/bin/bash
set -e

######## CREATE ########
# 프로젝트 폴더 생성
function create_project_directory() {
	if [ ! -d "/project" ]; then
		mkdir /project
	fi
}

# create & clone docker
function create_docker() {
	if [ ! -d "/dockers" ]; then
		echo "Clone Dockers..."
		pushd /
		git clone https://github.com/0strich/dockers.git
		popd
		echo "Dockers Cloned"
	fi
}

######## CLONE ########
# nginx clone
function clone_nginx() {
	create_project_directory
	create_docker
	if [ ! -d "/project/nginx" ]; then
		cd /project
		git clone https://github.com/0strich/nginx.git
		cd /project/nginx
		source start.sh
	fi
}

# mongodb clone
function clone_mongodb() {
	create_project_directory
	create_docker
	if [ ! -d "/project/mongodb" ]; then
		cp -r /dockers/mongodb /project/mongodb
	fi
}
