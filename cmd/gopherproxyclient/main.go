package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/CanadianCommander/gopherproxy/cmd/gopherproxyclient/proxy"
	"github.com/CanadianCommander/gopherproxy/internal/logging"
	"github.com/CanadianCommander/gopherproxy/internal/proxcom"
	proxylib "github.com/CanadianCommander/gopherproxy/internal/proxy"
	"go.uber.org/zap"
)

func main() {
	cliArgs := ParseArgs()
	if cliArgs.Debug {
		logging.CreateLogger(zap.DebugLevel)
	} else {
		logging.CreateLogger(zap.ErrorLevel)
	}

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
	clientManager.WaitForInitialization()

	switch cliArgs.Command {
	case "list":
		listChannelMembers(cliArgs.Channel, clientManager)
	case "echo":
		printLoop(client)
	}

	if !client.Closed {
		client.Close()
	}
}

func listChannelMembers(channel string, clientManager *proxy.ClientManager) {
	fmt.Printf("================ Clients On Channel [%s] ================\n", channel)
	for _, member := range clientManager.StateManager.ChannelMembers {
		if member.Id == clientManager.Client.Id {
			fmt.Printf("  %s (You)\n", member.Name)
		} else {
			fmt.Printf("  %s \n", member.Name)
		}
	}
}

func printLoop(client *proxylib.ProxyClient) {

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
