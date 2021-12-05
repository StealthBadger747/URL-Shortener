#!/bin/bash

REDIS_IP=127.0.0.1

if [[ $(uname -a | grep Linux) ]];
then
    echo "Linux machine detected"
    REDIS_IP=$(ip -o -4  address show  | awk ' NR==2 { gsub(/\/.*/, "", $4); print $4 }')
else
    echo "Other *NIX system machine detected"
    REDIS_IP=$(ifconfig | grep "inet " | grep -Fv 127.0.0.1 | awk '{print $2}')
fi

echo $REDIS_IP

REDIS_IP=$REDIS_IP docker-compose down
docker build . -t url-shortener
REDIS_IP=$REDIS_IP docker-compose up