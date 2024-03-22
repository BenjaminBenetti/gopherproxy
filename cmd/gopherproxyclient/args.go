package main

import (
	"flag"
)

type CliArgs struct {
	ProxyUrl string
	Password string
}

// ============================================
// Public Methods
// ============================================

// ParseArgs parses the command line arguments
func ParseArgs() CliArgs {

	proxyUrl := flag.String("proxy", "wss://localhost", "The URL of the GopherProxy instance")
	password := flag.String("password", "", "The password to use for the proxy connection")
	flag.Parse()

	cliArgs := CliArgs{
		ProxyUrl: *proxyUrl,
		Password: *password,
	}

	validateArgs(cliArgs)
	return cliArgs
}

// ============================================
// Private Methods
// ============================================

func validateArgs(args CliArgs) {
	if args.ProxyUrl == "" {
		panic("You must provide a proxy URL to connect to. Type --help for more information.")
	}
	if args.Password == "" {
		panic("You must provide a password to connect to the proxy. Type --help for more information.")
	}
}
