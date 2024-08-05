# DevCluster

This directory contains scripts to setup a local kubernetes cluster for development.

## Setup
Run the setup script. This will create a local kubernetes cluster.
```bash
./setup.sh 
```
That's it.

## Teardown 
Run the teardown script to delete the local kubernetes cluster.
```bash
./teardown.sh
```
Cluster setup & teardown is so lightweight that you can do it as often as you like. Perhaps to keep the clutter down on 
your computer you might want to run the teardown script after each time you're done working.
