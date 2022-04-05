#!/bin/bash

. scripts/utils.sh

for channel in mamichannel mischannel micchannel scchannel cmachannel; do
    infoln "Creating channel: $channel"
    if [ "$channel" == "mamichannel" ]; then
        ORG0=manufacturer
        ORG0_MSP=ManufacturerMSP
        ORG1=middleman
        ORG1_MSP=MiddlemanMSP
    elif [ "$channel" == "mischannel" ]; then
        ORG0=middleman
        ORG0_MSP=MiddlemanMSP
        ORG1=supplier
        ORG1_MSP=SupplierMSP
    elif [ "$channel" == "micchannel" ]; then
        ORG0=middleman
        ORG0_MSP=MiddlemanMSP
        ORG1=carrier
        ORG1_MSP=CarrierMSP
    elif [ "$channel" == "scchannel" ]; then
        ORG0=supplier
        ORG0_MSP=SupplierMSP
        ORG1=carrier
        ORG1_MSP=CarrierMSP
    else
        ORG0=carrier
        ORG0_MSP=CarrierMSP
        ORG1=manufacturer
        ORG1_MSP=ManufacturerMSP
    fi
    export CORE_PEER_MSPCONFIGPATH=/tmp/hyperledger/fabric/peer/$ORG0/admin/msp
    export CORE_PEER_TLS_ROOTCERT_FILE=/tmp/hyperledger/fabric/peer/$ORG0/peer1/tls/tlscacerts/tls-0-0-0-0-7052.pem
    export CORE_PEER_TLS_CERT_FILE=/tmp/hyperledger/fabric/peer/$ORG0/peer1/tls/signcerts/cert.pem
    export CORE_PEER_TLS_KEY_FILE=/tmp/hyperledger/fabric/peer/$ORG0/peer1/tls/keystore/key.pem
    export CORE_PEER_LOCALMSPID=$ORG0_MSP

    peer channel create -c $channel -f /tmp/hyperledger/fabric/channel-artifacts/$channel.tx -o orderer-cbpm:7050 --outputBlock /tmp/hyperledger/fabric/channel-artifacts/$channel.block --tls --cafile /tmp/hyperledger/fabric/peer/cbpm/orderer/tls/tlscacerts/tls-0-0-0-0-7052.pem 
    if [ $? -ne 0 ]; then
        fatalln "fail to create channel: $channel"
    fi
    successln "$channel successfully created"

    export CORE_PEER_ADDRESS=peer1-$ORG0:7051
    peer channel join -b /tmp/hyperledger/fabric/channel-artifacts/$channel.block 
    if [ $? -ne 0 ]; then
        fatalln "$CORE_PEER_ADDRESS fail to join channel $channel"
    fi
    successln "$CORE_PEER_ADDRESS successfully join channel $channel "

    export CORE_PEER_ADDRESS=peer2-$ORG0:7051
    peer channel join -b /tmp/hyperledger/fabric/channel-artifacts/$channel.block 
    if [ $? -ne 0 ]; then
        fatalln "$CORE_PEER_ADDRESS fail to join channel $channel"
    fi
    successln "$CORE_PEER_ADDRESS successfully join channel $channel "

    export CORE_PEER_MSPCONFIGPATH=/tmp/hyperledger/fabric/peer/$ORG1/admin/msp
    export CORE_PEER_TLS_ROOTCERT_FILE=/tmp/hyperledger/fabric/peer/$ORG1/peer1/tls/tlscacerts/tls-0-0-0-0-7052.pem
    export CORE_PEER_TLS_CERT_FILE=/tmp/hyperledger/fabric/peer/$ORG1/peer1/tls/signcerts/cert.pem
    export CORE_PEER_TLS_KEY_FILE=/tmp/hyperledger/fabric/peer/$ORG1/peer1/tls/keystore/key.pem
    export CORE_PEER_LOCALMSPID=$ORG1_MSP

    export CORE_PEER_ADDRESS=peer1-$ORG1:7051
    peer channel join -b /tmp/hyperledger/fabric/channel-artifacts/$channel.block 
    if [ $? -ne 0 ]; then
        fatalln "$CORE_PEER_ADDRESS fail to join channel $channel"
    fi
    successln "$CORE_PEER_ADDRESS successfully join channel $channel "
    export CORE_PEER_ADDRESS=peer2-$ORG1:7051
    peer channel join -b /tmp/hyperledger/fabric/channel-artifacts/$channel.block 
    if [ $? -ne 0 ]; then
        fatalln "$CORE_PEER_ADDRESS fail to join channel $channel"
    fi
    successln "$CORE_PEER_ADDRESS successfully join channel $channel "
done
