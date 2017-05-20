#!/bin/sh
echo Building app..

cd producer
./build.sh
cd ..
cd consumer
./build.sh
cd ..

docker-compose -f docker-compose-build.yml up -d