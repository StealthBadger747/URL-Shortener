#!/bin/bash

# These are the configuration envrionment variables
SERVER_PORT=8089
FRONTEND_DIR="$PWD/FrontendHtmx"
USE_REDIS=FALSE
# Uncomment if USE_REDIS=TRUE
#REDIS_IP={MACHINE IP}
#REDIS_PORT={THE REDIS PORT}

# # Build the backend
cd Backend
mvn install
mvn package
cd ..

# Run the jar and pass the ENV variables
FRONTEND_DIR=$FRONTEND_DIR \
    SERVER_PORT=$SERVER_PORT \
    USE_REDIS=$USE_REDIS \
    java -jar ./Backend/target/Url-Shortener-1.0.jar
