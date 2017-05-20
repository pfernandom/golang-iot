#!/bin/sh
echo Building producer..

docker build -t producer:build . -f Dockerfile.build

docker create --name extract-producer producer:build  
docker cp extract-producer:/go/src/github.com/pfernandom/streams/main ./main  
docker rm -f extract-producer

echo Creating image

docker build --no-cache -t pfernandom/producer:latest . -f Dockerfile.run