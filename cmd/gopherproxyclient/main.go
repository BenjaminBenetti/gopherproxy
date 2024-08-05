package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/CanadianCommander/gopherproxy/cmd/gopherproxyclient/proxy"
	"github.com/CanadianCommander/gopherproxy/internal/proxcom"
	proxylib "github.com/CanadianCommander/gopherproxy/internal/proxy"
)

func main() {
	cliArgs := ParseArgs()

	// Create a new GopherProxyClient
	client, err := proxylib.NewOutgoingSocket(cliArgs.ProxyUrl, proxylib.ProxyClientSettings{
		Channel:  cliArgs.Channel,
		Password: cliArgs.Password,
		Name:     cliArgs.ClientName,
	})

	if err != nil {
		fmt.Print("Failed to connect to GopherProxy server")
		panic(err)
	}

	clientManager := proxy.NewClientManager(client)
	clientManager.Start()

	for {
		fmt.Print("> ")
		reader := bufio.NewReader(os.Stdin)
		command, _ := reader.ReadString('\n')

		packet := proxcom.Packet{
			Type: proxcom.Data,
			Target: proxcom.Endpoint{
				Ip:   "0.0.0.0",
				Port: 70,
				Name: "gopherserver",
			},
			Source: proxcom.Endpoint{
				Ip:   "127.0.0.1",
				Port: 12345,
				Name: "gopherproxyclient",
			},
			Data: []byte(command),
		}

		client.Write(packet)
	}
}
