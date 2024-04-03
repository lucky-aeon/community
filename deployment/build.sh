#!/bin/bash
cd $(dirname "$(readlink -f "$0")")
rm -rf ../build 

echo "Building the back end..."
(
    cd ../
    go mod tidy
    BUILD_PATH="./build/community"
    # if [ $OSTYPE == "msys"]
    # then
    #     BUILD_PATH="./build/community.exe"
    GOOS=linux go build -ldflags "-s -w" -o ${BUILD_PATH} ./cmd/community
    # cp -r ./cmd/community/config.yaml ./build/config.yaml
    cp ./deployment/start.sh ./build/start.sh
    echo "build finished"
)

echo "Building the front end..."
(
    cd ../web
    yarn
    yarn build
    cp -r ./dist ../build/web/
    echo "build finished"   
)
