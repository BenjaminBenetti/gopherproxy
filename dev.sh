#!/bin/bash

if [[ $1 == "--help" ]] || [[ $1 == "-h" ]]; then
  echo "Usage: ./dev.sh [--auto-clean]"
  echo "  --auto-clean: Automatically clean up the cluster after dev is done."
  exit 0
fi

# all in one dev script. Creates cluster, deploys, and tears down all in one command.
pushd $(dirname $0) >> /dev/null

export PROJECT_NAME=$(cat ./project.json | jq -r '.name')

./devcluster/setup.sh
skaffold dev

if [[ $1 == "--auto-clean" ]]; then
  ./devcluster/teardown.sh
fi

popd >> /dev/null