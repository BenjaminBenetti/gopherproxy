#!/bin/bash
pushd $(dirname $0) >> /dev/null

minikube delete --all

popd >> /dev/null