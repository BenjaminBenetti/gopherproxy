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
- [DevContainers](https://containers.dev/) - An IDE that supports devcontainers. i.e. `vscode`

# Running
Simply open the project devcontainer and run the following command:

```bash
./dev.sh
```

Your proxy should now be running at `proxy.gopherproxy.dev`. You can now connect to the proxy server with the client.

```bash
go run ./cmd/gopherproxyclient/ --proxy 'wss://proxy.gopherproxy.dev/api/ws/connect' --password abc123 --channel test --name bobross
```