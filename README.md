# kube-dev-template

If you're anything like me, you've spent a lot of time setting up development environments for various projects.
You, of course, want to develop inside Kubernetes, overcomplicating things, just like me. :sweat_smile:
I've finally setup a the same scripts/helm/skaffold files one to many times. So I've made this template which should get you up and running 
quickly :racehorse:. I think you can be up and running in `10 - 20` minutes with this template. 

This is a template for a Kubernetes native project. Including local k8s dev setup as well as prod deployment manifests. 
This project let's you quickly bootstrap your kubernetes development environment! :rocket: 

## Prerequisites

### Things you probably already have 
- [Linux](https://www.linux.org/) - only compatible with Linux... May work with [WSL](https://learn.microsoft.com/en-us/windows/wsl/install) but I haven't tested it because Windows is gross. :poop:
- [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl-linux/) - to interact with k8s
- [kind](https://kind.sigs.k8s.io/docs/user/quick-start/) - to setup a local k8s cluster on your computer
- [helm](https://helm.sh/docs/intro/install/) - to deploy manifests to k8s

### Things you might not
- [skaffold](https://skaffold.dev/docs/install/) - to develop in your cluster 
- [mkcert](https://github.com/FiloSottile/mkcert) - to create local development certificates, so that SSL works!

## Setup

Once you've installed the prerequisites, edit the following files to configure your project:
1. `project.json` - This file contains general information about your project. For example, the name of your project.
2. `helm/services.yaml` - This file contains the project specific list of services (microservices) that make up your project. This file defines how to deploy your project to kubernetes. Following the provided example in that file add your own services.
3. `skaffold.yaml` - Add build definitions for project specific images (dockerfiles) to the artifacts list. 

## Usage

Run your app locally in dev mode with the following command. :rocket:

```bash 
./dev.sh [--auto-clean]
```
