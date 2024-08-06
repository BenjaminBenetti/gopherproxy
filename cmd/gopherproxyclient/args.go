package main

import (
	"flag"
	"fmt"
	"net/url"
	"os"

	"github.com/CanadianCommander/gopherproxy/cmd/gopherproxyclient/proxy"
)

type CliArgs struct {
	ProxyUrl        url.URL
	Password        string
	Channel         string
	ClientName      string
	Debug           bool
	Command         string
	ForwardingRules []*proxy.ForwardingRule
}

// ============================================
// Public Methods
// ============================================

// ParseArgs parses the command line arguments
func ParseArgs() CliArgs {
	setupHelpMessage()

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

	// parse fowarding rules
	forwardingRules := make([]*proxy.ForwardingRule, 0)
	for i := 1; i < flag.NArg(); i++ {
		forwardingRules = append(forwardingRules, proxy.NewForwardingRuleFromArg(flag.Arg(i)))
	}

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
		ProxyUrl:        *proxyUrl,
		Password:        *password,
		Channel:         *channel,
		ClientName:      *clientName,
		Debug:           *debug,
		Command:         command,
		ForwardingRules: forwardingRules,
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

func setupHelpMessage() {
	flag.Usage = func() {
		_, _ = os.Stderr.WriteString("Usage: gopherproxyclient [options] <command> [<forward definition>]\n")
		fmt.Println("Commands:")
		fmt.Println("  list    List all clients connected to the channel")
		fmt.Println("  start   Start the client and forward traffic as defined by the forward definitions")
		fmt.Println("Forward Defninition:")
		fmt.Println("  <forward definition> defines how traffic should be proxied. It can appear multiple times. It has the following format:")
		fmt.Println("  <local port>:<remote client>:[remote host]:<remote port>")
		fmt.Println("    - local port: The port on the local machine to listen on")
		fmt.Println("    - remote client: The name of the remote client to forward traffic to. Use the \"list\" command to see available clients.")
		fmt.Println("    - [optional] remote host: The host to forward traffic to on the remote client. Defaults to localhost.")
		fmt.Println("    - remote port: The port to forward traffic to on the remote target.")
		fmt.Println("  Example: 8080:client1:google.com:80 - Forward traffic on local port 8080 to google.com:80 from client1")
		fmt.Println("Options:")
		flag.PrintDefaults()
	}
}
