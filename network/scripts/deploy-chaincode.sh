#!/bin/bash

. scripts/utils.sh

channel=cbpmchannel
CHAINCODE_LANG=golang
CHAINCODE_VERSION=1.0
CHAINCODE=chaincode-cbpm
CHAINCODE_NAME=cbpmchaincode

for ORG in manufacturer middleman carrier supplier; do
    if [ "$ORG" == "manufacturer" ]; then
        ORG_MSP=ManufacturerMSP
    elif [ "$ORG" == "middleman" ]; then
        ORG_MSP=MiddlemanMSP
    elif [ "$ORG" == "supplier" ]; then
        ORG_MSP=SupplierMSP
    else
        ORG_MSP=CarrierMSP
    fi

    export CORE_PEER_TLS_ROOTCERT_FILE=/tmp/hyperledger/fabric/peer/$ORG/peer1/tls/tlscacerts/tls-0-0-0-0-7052.pem
    export CORE_PEER_TLS_CERT_FILE=/tmp/hyperledger/fabric/peer/$ORG/peer1/tls/signcerts/cert.pem
    export CORE_PEER_TLS_KEY_FILE=/tmp/hyperledger/fabric/peer/$ORG/peer1/tls/keystore/key.pem
    export CORE_PEER_MSPCONFIGPATH=/tmp/hyperledger/fabric/peer/$ORG/admin/msp
    export CORE_PEER_ADDRESS=peer1-$ORG:7051
    export CORE_PEER_LOCALMSPID=$ORG_MSP

    infoln "Installing chaincode $CHAINCODE_NAME on $CORE_PEER_ADDRESS"    

    peer chaincode install  -n $CHAINCODE_NAME -v $CHAINCODE_VERSION -p $CHAINCODE -l $CHAINCODE_LANG 

    if [ $? -ne 0 ]; then
        fatalln "fail to install chaincode $CHAINCODE_NAME on $CORE_PEER_ADDRESS"
    fi
    successln "chaincode successfully installed $CHAINCODE_NAME on $CORE_PEER_ADDRESS"

    export CORE_PEER_ADDRESS=peer2-$ORG:7051

    infoln "Installing chaincode $CHAINCODE_NAME on $CORE_PEER_ADDRESS"
     
    peer chaincode install  -n $CHAINCODE_NAME -v $CHAINCODE_VERSION -p $CHAINCODE -l $CHAINCODE_LANG 

    if [ $? -ne 0 ]; then
        fatalln "fail to install chaincode $CHAINCODE_NAME on $CORE_PEER_ADDRESS"
    fi
    successln "chaincode successfully installed $CHAINCODE_NAME on $CORE_PEER_ADDRESS"
done


infoln "Instantiating chaincode $CHAINCODE_NAME"  
peer chaincode instantiate  -o orderer-cbpm:7050 --tls --cafile /tmp/hyperledger/fabric/peer/cbpm/orderer/tls/tlscacerts/tls-0-0-0-0-7052.pem -C $channel -n $CHAINCODE_NAME -l $CHAINCODE_LANG -v $CHAINCODE_VERSION -c '{"Args":[""]}' -P "OR('ManufacturerMSP.peer','MiddlemanMSP.peer','SupplierMSP.peer','CarrierMSP.peer')" --collections-config $GOPATH/src/$CHAINCODE/collections_config.json
if [ $? -ne 0 ]; then
    fatalln "fail to instantiate chaincode $CHAINCODE_NAME"
fi
successln "chaincode successfully instantiated $CHAINCODE_NAME"
