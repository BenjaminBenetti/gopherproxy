# DevCluster

This directory contains scripts to setup a local kubernetes cluster for development.

## Prerequisites
- [kind](https://kind.sigs.k8s.io/docs/user/quick-start/)
- [helm](https://helm.sh/docs/intro/install/)
- [skaffold](https://skaffold.dev/docs/install/)

## Setup
Run the setup script. This will create a local kubernetes cluster.
```bash
./setup.sh 
```
Next edit your `/etc/hosts` file to include the following entries:
```
172.17.0.1 dev.devwiki.io
```
Finally run `skaffold dev` from the root of this project to boot devwiki in dev mode. 
You should now be able to see devwiki running in your browser at [dev.devwiki.io](http://dev.devwiki.io).

## Teardown 
Run the teardown script to delete the local kubernetes cluster.
```bash
./teardown.sh
```
Cluster setup & teardown is so lightweight that you can do it as often as you like. Perhaps to keep the clutter down on 
your computer you might want to run the teardown script after each time you're done working on devwiki. Then run the
setup script again next time you want to work on devwiki.
