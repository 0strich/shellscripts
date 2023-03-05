#!/bin/bash

# 설명
function printHelp() {
  echo 'print help'
}

# [v2] 패브릭 네트워크 실행
function fabricUp() {
  # download 하이퍼레저 패브릭 & docker pull
  curl -sSL https://bit.ly/2ysbOFE | bash -s -- 2.2.2 1.4.9

  # 폴더 이동
  pushd fabric-samples/test-network

  # 테스트넷 실행
  ./network.sh up createChannel
  ./network.sh deployCC -ccn basic -ccp ../asset-transfer-basic/chaincode-javascript -ccl javascript

  # 환경설정
  export PATH=${PWD}/../bin:$PATH
  export FABRIC_CFG_PATH=$PWD/../config/
  # Environment variables for Org1
  export CORE_PEER_TLS_ENABLED=true
  export CORE_PEER_LOCALMSPID="Org1MSP"
  export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt
  export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp
  export CORE_PEER_ADDRESS=localhost:7051

  # 피어 초기화
  peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile ${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C mychannel -n basic --peerAddresses localhost:7051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt --peerAddresses localhost:9051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt -c '{"function":"InitLedger","Args":[]}'

  # 조회
  peer chaincode query -C mychannel -n basic -c '{"Args":["GetAllAssets"]}'

  popd
}

if [ "$1" == "up" ]; then
  fabricUp
else
  printHelp
  exit 1
fi
