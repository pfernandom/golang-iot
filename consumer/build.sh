#!/bin/sh
echo Building consumer..

docker build -t consumer:build . -f Dockerfile.build

docker create --name extract-consumer consumer:build  
docker cp extract-consumer:/go/src/github.com/pfernandom/streams/main ./main  
docker rm -f extract-consumer

echo Creating image

docker build --no-cache -t pfernandom/consumer:latest . -f Dockerfile.run