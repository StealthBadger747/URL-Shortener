#!/bin/bash

# These are the configuration envrionment variables
SERVER_PORT=8089
ANGULAR_FRONTEND_DIR="$PWD/Frontend/dist/UrlShortener/"
USE_REDIS=FALSE

# Build the frontend
cd Frontend
npm install
npx ng build
cd ..

# # Build the backend
cd backend
mvn install
mvn package
cd ..

# Run the jar and pass the ENV variables
ANGULAR_FRONTEND_DIR=$ANGULAR_FRONTEND_DIR \
    SERVER_PORT=$SERVER_PORT \
    USE_REDIS=$USE_REDIS \
    java -jar ./Backend/target/Url-Shortener-1.0.jar
