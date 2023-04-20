#!/bin/bash
set -e

# apt 업데이트 & vim 설치
apt update -y && apt install vim net-tools -y

mkdir /project && cd /project

# go설치 & 압축 해제 & 경로 설정 & 적용
wget https://golang.org/dl/go1.15.2.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.15.2.linux-amd64.tar.gz
rm go1.15.2.linux-amd64.tar.gz
echo "export PATH=$PATH:/usr/local/go/bin" >>~/.bashrc
source ~/.bashrc

# go-ipfs 설치 & ipfs 명령어 폴더 이동 & 경로 설정
wget https://dist.ipfs.io/go-ipfs/v0.6.0/go-ipfs_v0.6.0_linux-amd64.tar.gz
tar zxvf go-ipfs_v0.6.0_linux-amd64.tar.gz
rm go-ipfs_v0.6.0_linux-amd64.tar.gz
mv go-ipfs/ipfs /usr/bin/ipfs
rm -rf ./go-ipfs
echo "IPFS_PATH=~/.ipfs" >>~/.bashrc
source ~/.bashrc

# ipfs 초기화
ipfs init
ipfs daemon &
