#!/bin/bash

export PATH=${PWD}/../../bin:${PWD}:$PATH
export FABRIC_CFG_PATH=$PWD/../../config/
export FABRIC_CA_CLIENT_HOME=${PWD}/../organizations/peerOrganizations/org1.example.com/

fabric-ca-client register --caname ca-org1 --id.name $1 --id.secret $2 --id.type client --tls.certfiles ${PWD}/../organizations/fabric-ca/org1/tls-cert.pem

fabric-ca-client enroll -u https://$1:$2@localhost:7054 --caname ca-org1 -M ${PWD}/../organizations/peerOrganizations/org1.example.com/users/$1@org1.example.com/msp --tls.certfiles ${PWD}/../organizations/fabric-ca/org1/tls-cert.pem
cp ${PWD}/../organizations/peerOrganizations/org1.example.com/msp/config.yaml ${PWD}/../organizations/peerOrganizations/org1.example.com/users/$1@org1.example.com/msp/config.yaml
