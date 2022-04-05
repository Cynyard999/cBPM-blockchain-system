#!/bin/bash

. scripts/utils.sh

CONTAINER_CLI="docker"
CONTAINER_CLI_COMPOSE="${CONTAINER_CLI}-compose"

infoln "Using ${CONTAINER_CLI} and ${CONTAINER_CLI_COMPOSE}"

infoln "Shutting down network"
${CONTAINER_CLI_COMPOSE} down


infoln "Removing remaining containers"
${CONTAINER_CLI} rm -f $(${CONTAINER_CLI} ps -aq --filter name='dev-peer*') 2>/dev/null || true
infoln "Removing networks"
${CONTAINER_CLI} network prune -f
infoln "Removing volumes"
${CONTAINER_CLI} volume prune -f
