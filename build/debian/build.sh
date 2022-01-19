#!/usr/bin/env bash

#sudo apt-get update
#sudo apt-get install -y ruby-dev build-essential
#sudo gem install fpm

rm -rf ../../target/debian
mkdir -p ../../target/debian/package/usr/bin
cd ../../target/debian/package/usr/bin
go build ../../../../../cmd/docker-build
go build ../../../../../cmd/docker-devcontainer
mkdir -p ../share/docker-tools
cp -R ../../../../../templates ../share/docker-tools 

PKG_NAME=docker-tools
PKG_DESCRIPTION="Extended tools for docker"
PKG_VERSION=1.1.0
PKG_RELEASE=1
PKG_MAINTAINER="Brian Bintz <brian@bintzpress.com>"
PKG_VENDOR="Bintz Press"
PKG_URL="https://bintzpress.com"

FPM_OPTS="-n $PKG_NAME -v $PKG_VERSION --iteration $PKG_RELEASE"

cd ../../../
CUR_DIR=`pwd`
echo "at ".$CUR_DIR

fpm -s dir -t deb ${FPM_OPTS} -f \
    --maintainer "$PKG_MAINTAINER" \
    --vendor "$PKG_VENDOR" \
    --url "$PKG_URL" \
    --description "$PKG_DESCRIPTION" \
    --architecture "amd64" \
    -C "package" \
    usr 

