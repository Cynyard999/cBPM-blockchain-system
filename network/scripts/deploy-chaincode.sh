#!/bin/bash

. scripts/utils.sh

for channel in mamichannel mischannel micchannel scchannel cmachannel; do
    infoln "Deploying chaincode on $channel"
    if [ "$channel" == "mamichannel" ]; then
        ORG0=manufacturer
        ORG0_MSP=ManufacturerMSP
        ORG1=middleman
        ORG1_MSP=MiddlemanMSP
        CHAINCODE=chaincode-manufacturer-middleman
        CHAINCODE_NAME=mamichaincode
    elif [ "$channel" == "mischannel" ]; then
        ORG0=middleman
        ORG0_MSP=MiddlemanMSP
        ORG1=supplier
        ORG1_MSP=SupplierMSP
        CHAINCODE=chaincode-middleman-supplier
        CHAINCODE_NAME=mischaincode
    elif [ "$channel" == "micchannel" ]; then
        ORG0=middleman
        ORG0_MSP=MiddlemanMSP
        ORG1=carrier
        ORG1_MSP=CarrierMSP
        CHAINCODE=chaincode-middleman-carrier
        CHAINCODE_NAME=micchaincode
    elif [ "$channel" == "scchannel" ]; then
        ORG0=supplier
        ORG0_MSP=SupplierMSP
        ORG1=carrier
        ORG1_MSP=CarrierMSP
        CHAINCODE=chaincode-supplier-carrier
        CHAINCODE_NAME=scchaincode
    else
        ORG0=carrier
        ORG0_MSP=CarrierMSP
        ORG1=manufacturer
        ORG1_MSP=ManufacturerMSP
        CHAINCODE=chaincode-carrier-manufacturer
        CHAINCODE_NAME=cmachaincode
    fi

    CHAINCODE_LANG=golang
    CHAINCODE_VERSION=1.0

    export CORE_PEER_TLS_ROOTCERT_FILE=/tmp/hyperledger/fabric/peer/$ORG0/peer1/tls/tlscacerts/tls-0-0-0-0-7052.pem
    export CORE_PEER_TLS_CERT_FILE=/tmp/hyperledger/fabric/peer/$ORG0/peer1/tls/signcerts/cert.pem
    export CORE_PEER_TLS_KEY_FILE=/tmp/hyperledger/fabric/peer/$ORG0/peer1/tls/keystore/key.pem
    export CORE_PEER_MSPCONFIGPATH=/tmp/hyperledger/fabric/peer/$ORG0/admin/msp
    export CORE_PEER_ADDRESS=peer1-$ORG0:7051
    export CORE_PEER_LOCALMSPID=$ORG0_MSP

    infoln "Installing chaincode $CHAINCODE_NAME on $CORE_PEER_ADDRESS"    

    peer chaincode install  -n $CHAINCODE_NAME -v $CHAINCODE_VERSION -p $CHAINCODE -l $CHAINCODE_LANG 

    if [ $? -ne 0 ]; then
        fatalln "fail to install chaincode $CHAINCODE_NAME on $CORE_PEER_ADDRESS"
    fi
    successln "chaincode successfully installed $CHAINCODE_NAME on $CORE_PEER_ADDRESS"

    export CORE_PEER_ADDRESS=peer2-$ORG0:7051

    infoln "Installing chaincode $CHAINCODE_NAME on $CORE_PEER_ADDRESS"  
    peer chaincode install  -n $CHAINCODE_NAME -v $CHAINCODE_VERSION -p $CHAINCODE -l $CHAINCODE_LANG 

    if [ $? -ne 0 ]; then
        fatalln "fail to install chaincode $CHAINCODE_NAME on $CORE_PEER_ADDRESS"
    fi
    successln "chaincode successfully installed $CHAINCODE_NAME on $CORE_PEER_ADDRESS"


    export CORE_PEER_MSPCONFIGPATH=/tmp/hyperledger/fabric/peer/$ORG1/admin/msp
    export CORE_PEER_ADDRESS=peer1-$ORG1:7051
    export CORE_PEER_LOCALMSPID=$ORG1_MSP
    export CORE_PEER_TLS_ROOTCERT_FILE=/tmp/hyperledger/fabric/peer/$ORG1/peer1/tls/tlscacerts/tls-0-0-0-0-7052.pem
    export CORE_PEER_TLS_CERT_FILE=/tmp/hyperledger/fabric/peer/$ORG1/peer1/tls/signcerts/cert.pem
    export CORE_PEER_TLS_KEY_FILE=/tmp/hyperledger/fabric/peer/$ORG1/peer1/tls/keystore/key.pem


    infoln "Installing chaincode $CHAINCODE_NAME on $CORE_PEER_ADDRESS"  
    peer chaincode install  -n $CHAINCODE_NAME -v $CHAINCODE_VERSION -p $CHAINCODE -l $CHAINCODE_LANG 
    if [ $? -ne 0 ]; then
        fatalln "fail to install chaincode $CHAINCODE_NAME on $CORE_PEER_ADDRESS"
    fi
    successln "chaincode successfully installed $CHAINCODE_NAME on $CORE_PEER_ADDRESS"


    export CORE_PEER_ADDRESS=peer2-$ORG1:7051

    infoln "Installing chaincode $CHAINCODE_NAME on $CORE_PEER_ADDRESS"  
    peer chaincode install  -n $CHAINCODE_NAME -v $CHAINCODE_VERSION -p $CHAINCODE -l $CHAINCODE_LANG 
    if [ $? -ne 0 ]; then
        fatalln "fail to install chaincode $CHAINCODE_NAME on $CORE_PEER_ADDRESS"
    fi
    successln "chaincode successfully installed $CHAINCODE_NAME on $CORE_PEER_ADDRESS"


    infoln "Instantiating chaincode $CHAINCODE_NAME"  
    peer chaincode instantiate  -o orderer-cbpm:7050 --tls --cafile /tmp/hyperledger/fabric/peer/cbpm/orderer/tls/tlscacerts/tls-0-0-0-0-7052.pem -C $channel -n $CHAINCODE_NAME -l $CHAINCODE_LANG -v $CHAINCODE_VERSION -c '{"Args":[""]}' -P "OR('$ORG0_MSP.peer','$ORG1_MSP.peer')" 
    if [ $? -ne 0 ]; then
        fatalln "fail to instantiate chaincode $CHAINCODE_NAME"
    fi
    successln "chaincode successfully instantiated $CHAINCODE_NAME"
done
