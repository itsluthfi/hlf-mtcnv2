#!/bin/bash

./network.sh down

./network.sh up createChannel -s couchdb -c mychannel -ca

# ./network.sh deployCC -ccn ledger -ccp ../chaincode-go/ -ccl go -ccep "OR('Org1MSP.peer','Org2MSP.peer')"

# ./network.sh deployCC -ccn ledger -ccp ../chaincode-go/ -ccl go -ccep "AND('Org1MSP.peer','Org2MSP.peer')"

./network.sh deployCC -ccn basic -ccp ../token-erc-20/chaincode-go/ -ccl go -ccep "AND('Org1MSP.peer','Org2MSP.peer')"
