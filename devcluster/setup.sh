#!/bin/bash

pushd $(dirname $0) >> /dev/null

PROJECT_NAME=$(cat ../project.json | jq -r '.name')
PROJECT_DNS=${PROJECT_NAME}.dev

# create minikube cluster
minikube start --extra-config=apiserver.service-node-port-range=1-65535
# update /etc/hosts file
echo "Exposing proxy.gopherproxy.dev on https://$(minikube ip):443"
sudo sed -i "/proxy.gopherproxy.dev/d" /etc/hosts
sudo bash -c "echo \"$(minikube ip) proxy.gopherproxy.dev\" >> /etc/hosts"

# install helm repos 
helm repo add traefik https://traefik.github.io/charts
helm repo add metrics-server https://kubernetes-sigs.github.io/metrics-server/
helm repo update 

# setup dev cert
mkcert -install
pushd /tmp/ >> /dev/null
mkcert *.${PROJECT_DNS}
popd >> /dev/null
kubectl create namespace traefik-ingress
kubectl create secret generic cluster-cert \
    --from-file=tls.crt=/tmp/_wildcard.${PROJECT_DNS}.pem \
    --from-file=tls.key=/tmp/_wildcard.${PROJECT_DNS}-key.pem --namespace traefik-ingress
kubectl apply -f ./manifest/traefik-config.yaml  -n traefik-ingress

# install traefik 
helm install traefik traefik/traefik --create-namespace --namespace traefik-ingress --values ./manifest/traefik-values.yaml

# install metrics server
helm upgrade --install metrics-server metrics-server/metrics-server --namespace kube-system --set args[0]="--kubelet-insecure-tls=true"

echo Cluster setup complete.

popd >> /dev/null