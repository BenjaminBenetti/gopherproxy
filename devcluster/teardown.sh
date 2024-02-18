#!/bin/bash


pushd $(dirname $0) >> /dev/null

PROJECT_NAME=$(cat ../project.json | jq -r '.name')

kind delete cluster --name ${PROJECT_NAME}-cluster

popd >> /dev/null