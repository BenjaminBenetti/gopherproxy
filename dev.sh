#!/bin/bash
pushd $(dirname $0) >> /dev/null

export PROJECT_NAME=$(cat ./project.json | jq -r '.name')

skaffold dev

popd >> /dev/null