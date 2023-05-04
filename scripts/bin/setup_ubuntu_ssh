#!/bin/bash
set -e

# root 로그인 허용
command sudo sed 's/^.*ssh-rsa/ssh-rsa/g' /root/.ssh/authorized_keys >authorized_keys
command sudo cp ./authorized_keys /root/.ssh/authorized_keys

file_path="/etc/ssh/sshd_config"

search_string1="LoginGraceTime"
search_string2="PermitRootLogin"
search_string3="StrictModes"

update_string1="LoginGraceTime yes"
update_string2="PermitRootLogin 120"
update_string3="StrictModes yes"

# ssh 설정
command sudo sed -i "s/^.*$search_string1.*$/$update_string1/g" $file_path
command sudo sed -i "s/^.*$search_string2.*$/$update_string2/g" $file_path
command sudo sed -i "s/^.*$search_string3.*$/$update_string3/g" $file_path
