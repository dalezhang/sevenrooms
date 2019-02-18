#!/bin/sh
cd $PROJECT_DIR/dist
cp -rf $PROJECT_DIR/Dockerfile .
cp -rf $PROJECT_DIR/config .
docker build -t $DOCKER_IMAGE .
docker push $DOCKER_IMAGE
