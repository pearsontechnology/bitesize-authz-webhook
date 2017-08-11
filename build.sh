#!/bin/bash

set -e

VERSION=${1:-0.0.1}
DOCKER_REPO="pearsontechnology/bitesize-authz-webhook"
DOCKER_IMAGE=${DOCKER_REPO}:${VERSION}

# run tests before build
echo "Running go test..."
go test

echo "Building application..."
git tag RELEASE-${VERSION}
git push origin RELEASE-${VERSION}
docker run --rm -v "$PWD":/usr/src/myapp -w /usr/src/myapp golang:1.6 bash -c "go get -d && CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bitesize-authz-webhook -v"

echo "Building docker image..."
docker build -t ${DOCKER_IMAGE} .
docker push ${DOCKER_IMAGE}
