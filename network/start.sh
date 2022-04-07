#!/bin/bash

. scripts/utils.sh

CONTAINER_CLI="docker"
CONTAINER_CLI_COMPOSE="${CONTAINER_CLI}-compose"

infoln "Using ${CONTAINER_CLI} and ${CONTAINER_CLI_COMPOSE}"

infoln "Starting network..."
${CONTAINER_CLI_COMPOSE} up -d

if [ $? -ne 0 ]; then
    fatalln "Fail to build up fabric network"
fi

successln "Successfully build up fabric network"

infoln "Initiating channels..."

${CONTAINER_CLI} exec -it cli /bin/bash -c "./scripts/init-channel.sh && ./scripts/deploy-chaincode.sh"