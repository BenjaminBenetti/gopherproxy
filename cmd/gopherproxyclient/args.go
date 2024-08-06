package main

import (
	"flag"
	"fmt"
	"net/url"
	"os"
)

type CliArgs struct {
	ProxyUrl   url.URL
	Password   string
	Channel    string
	ClientName string
	Debug      bool
	Command    string
}

// ============================================
// Public Methods
// ============================================

// ParseArgs parses the command line arguments
func ParseArgs() CliArgs {
	flag.Usage = func() {
		_, _ = os.Stderr.WriteString("Usage: gopherproxyclient [options] <command>\n")
		fmt.Println("Commands:")
		fmt.Println("  list    List all clients connected to the channel")
		fmt.Println("  echo    Simple testing mode where you can interactivly send messages to the channel")
		fmt.Println("Options:")
		flag.PrintDefaults()
	}

	proxyUrlStr := flag.String("proxy", "wss://localhost", "The URL of the GopherProxy instance")
	password := flag.String("password", "", "The password to use for the proxy connection")
	channel := flag.String("channel", "", "The channel to connect to. Use the same channel name on both ends of the connection.")
	clientName := flag.String("name", "", "The name of the client connecting to the proxy. Use this to organize clients. Defaults to the hostname of the machine.")
	debug := flag.Bool("debug", false, "Enable debug logging")

	flag.Parse()
	if flag.NArg() == 0 {
		flag.Usage()
		os.Exit(1)
	}

	// read first positional argument as command
	command := flag.Arg(0)

	proxyUrl, err := url.Parse(*proxyUrlStr)
	if err != nil {
		panic("The url provided for --proxy could not be parsed. Please provide a valid URL. --help for more information.")
	}

	if *clientName == "" {
		hostname, err := os.Hostname()
		if err != nil {
			*clientName = "unknown"
		} else {
			*clientName = hostname
		}
	}

	cliArgs := CliArgs{
		ProxyUrl:   *proxyUrl,
		Password:   *password,
		Channel:    *channel,
		ClientName: *clientName,
		Debug:      *debug,
		Command:    command,
	}

	validateArgs(cliArgs)
	return cliArgs
}

// ============================================
// Private Methods
// ============================================

func validateArgs(args CliArgs) {
	if args.Password == "" {
		panic("You must provide a password to connect to the proxy. Type --help for more information.")
	}
	if args.Channel == "" {
		panic("You must provide a channel to connect to. Type --help for more information.")
	}
}
