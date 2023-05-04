#!/bin/bash
set -e

USER_NAME=hyper

function fabric_init() {

	setup -d -n -g

	# 도커 User Add
	adduser $USER_NAME

	usermod -aG docker $USER_NAME

	chmod 666 /run/docker.sock

	sed -i "s/\(^sudo.*$\)/&$USER_NAME/" /etc/group

	service docker restart

	function hyper_commands() {
		git clone https://github.com/0strich/shellscripts.git
		source $HOME/shellscripts/init
	}

	# $USER_NAME 계정 설정
	su - $USER_NAME -c "
	$(declare -f hyper_commands)
	hyper_commands
"
}
