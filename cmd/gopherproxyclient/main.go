package main

import (
	"fmt"
	"os"

	"github.com/CanadianCommander/gopherproxy/cmd/gopherproxyclient/forwarddisplay"
	"github.com/CanadianCommander/gopherproxy/cmd/gopherproxyclient/proxy"
	"github.com/CanadianCommander/gopherproxy/internal/logging"
	proxylib "github.com/CanadianCommander/gopherproxy/internal/proxy"
	"go.uber.org/zap"
)

func main() {
	cliArgs := ParseArgs()
	if cliArgs.Debug {
		logging.CreateLogger(zap.DebugLevel)
	} else if cliArgs.LoggingBasedUi {
		logging.CreateLogger(zap.InfoLevel)
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

	clientManager := proxy.NewClientManager(client, cliArgs.ForwardingRules, cliArgs.ProxyUrl, cliArgs.DebugPrintPackets)
	clientManager.Start()
	clientManager.WaitForInitialization()

	switch cliArgs.Command {
	case "list":
		listChannelMembers(cliArgs.Channel, clientManager)
	case "start":
		var display forwarddisplay.Display
		if cliArgs.LoggingBasedUi {
			display = forwarddisplay.NewForwardLoggingUi(clientManager)
		} else {
			display = forwarddisplay.NewForwardUi(clientManager)
		}
		display.Build()
		if !cliArgs.Debug {
			display.Start()
		} else {
			// don't draw display in debug mode so you can see log output
			<-make(chan os.Signal, 1)
			logging.Get().Info("User requested exit")
			return
		}
	default:
		fmt.Printf("Unknown command: %s\n", cliArgs.Command)
	}

	clientManager.Close()
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
