# GopherProxy 

Simple two way TCP proxy over WebSockets. This proxy allows you 
to connect to remote systems that normally do not have direct connectivity between them. 

Why make this instead of using something off the shelf?
Because FUNNNNNN!

## Topology
The setup is pretty simple. Bot clients connect to 
the central proxy server, operating in an in-out mode. 
Then the proxy server simply relays the data between the two clients.
depending on the configuration of each client (what ports/addresses to forward)

# Prerequisites 
- [Linux](https://www.linux.org/) - only compatible with Linux... May work with [WSL](https://learn.microsoft.com/en-us/windows/wsl/install) but I haven't tested it because Windows is gross. :poop:
- [golang 1.21+](https://golang.org/doc/install) - to build the project locally (not strictly required, but nice to have)
- [docker](https://docs.docker.com/get-docker/) - to build the docker image
- [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl-linux/) - to interact with k8s
- [kind](https://kind.sigs.k8s.io/docs/user/quick-start/) - to setup a local k8s cluster on your computer
- [helm](https://helm.sh/docs/intro/install/) - to deploy manifests to k8s
- [skaffold](https://skaffold.dev/docs/install/) - to develop in your cluster 
- [mkcert](https://github.com/FiloSottile/mkcert) - to create local development certificates, so that SSL works!

# Running
Launch the app & development cluster with the following command:

```bash
./dev.sh
```

