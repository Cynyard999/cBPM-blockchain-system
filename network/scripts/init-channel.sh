#!/bin/bash

. scripts/utils.sh

channel=cbpmchannel
is_channel_created=0
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
    export CORE_PEER_MSPCONFIGPATH=/tmp/hyperledger/fabric/peer/$ORG/admin/msp
    export CORE_PEER_TLS_ROOTCERT_FILE=/tmp/hyperledger/fabric/peer/$ORG/peer1/tls/tlscacerts/tls-0-0-0-0-7052.pem
    export CORE_PEER_TLS_CERT_FILE=/tmp/hyperledger/fabric/peer/$ORG/peer1/tls/signcerts/cert.pem
    export CORE_PEER_TLS_KEY_FILE=/tmp/hyperledger/fabric/peer/$ORG/peer1/tls/keystore/key.pem
    export CORE_PEER_LOCALMSPID=$ORG_MSP

    if [ $is_channel_created -ne 1 ]; then
        infoln "Creating channel: $channel"
        peer channel create -c $channel -f /tmp/hyperledger/fabric/channel-artifacts/$channel.tx -o orderer-cbpm:7050 --outputBlock /tmp/hyperledger/fabric/channel-artifacts/$channel.block --tls --cafile /tmp/hyperledger/fabric/peer/cbpm/orderer/tls/tlscacerts/tls-0-0-0-0-7052.pem 
        if [ $? -ne 0 ]; then
            fatalln "fail to create channel: $channel"
        fi
        successln "$channel successfully created"
        is_channel_created=1
    fi

    export CORE_PEER_ADDRESS=peer1-$ORG:7051
    peer channel join -b /tmp/hyperledger/fabric/channel-artifacts/$channel.block 
    if [ $? -ne 0 ]; then
        fatalln "$CORE_PEER_ADDRESS fail to join channel $channel"
    fi
    successln "$CORE_PEER_ADDRESS successfully join channel $channel "

    export CORE_PEER_ADDRESS=peer2-$ORG:7051
    peer channel join -b /tmp/hyperledger/fabric/channel-artifacts/$channel.block 
    if [ $? -ne 0 ]; then
        fatalln "$CORE_PEER_ADDRESS fail to join channel $channel"
    fi
    successln "$CORE_PEER_ADDRESS successfully join channel $channel "
done
