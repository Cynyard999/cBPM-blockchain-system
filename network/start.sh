#!/bin/bash

. scripts/utils.sh

CONTAINER_CLI="docker"
CONTAINER_CLI_COMPOSE="${CONTAINER_CLI}-compose"

infoln "Using ${CONTAINER_CLI} and ${CONTAINER_CLI_COMPOSE}"

infoln "Starting network"
${CONTAINER_CLI_COMPOSE} up -d