package main

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/CanadianCommander/gopherproxy/internal/proxcom"
	"github.com/CanadianCommander/gopherproxy/internal/websocket"
)

func main() {
	cliArgs := ParseArgs()

	// Create a new GopherProxyClient
	client, err := websocket.NewOutgoingSocket(cliArgs.ProxyUrl, websocket.ProxyClientSettings{
		Channel:  cliArgs.Channel,
		Password: cliArgs.Password,
		Name:     cliArgs.ClientName,
	})

	if err != nil {
		fmt.Print("Failed to connect to GopherProxy server")
		panic(err)
	}

	go printLoop(client)
	createSigtermHandler(client)

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

func createSigtermHandler(client *websocket.ProxyClient) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		client.Close()
		os.Exit(0)
	}()
}

func printLoop(client *websocket.ProxyClient) {
	for {
		packet, ok := client.Read()
		if !ok {
			fmt.Println("Connection closed")
			return
		} else {
			fmt.Println(string(packet.Data))
		}
	}
}
