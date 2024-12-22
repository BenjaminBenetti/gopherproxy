package proxcom

import (
	"strconv"
	"strings"
)

type ForwardingRule struct {
	LocalPort    int
	RemoteClient string
	RemoteHost   string
	RemotePort   int
	// if based on the current state of the channel this rule is Valid
	Valid bool
}

// ============================================
// Constructors
// ============================================

// NewForwardingRuleFromArg creates a new forwarding rule
func NewForwardingRuleFromArg(arg string) *ForwardingRule {
	argSlic := strings.Split(arg, ":")

	switch len(argSlic) {
	case 3:
		localPort, err := strconv.Atoi(argSlic[0])
		if err != nil {
			panic("The local port provided for the forwarding rule could not be parsed. Please provide a valid port number.")
		}

		remotePort, err := strconv.Atoi(argSlic[2])
		if err != nil {
			panic("The remote port provided for the forwarding rule could not be parsed. Please provide a valid port number.")
		}

		return &ForwardingRule{
			LocalPort:    localPort,
			RemoteClient: argSlic[1],
			RemoteHost:   "localhost",
			RemotePort:   remotePort,
			Valid:        false,
		}
	case 4:
		localPort, err := strconv.Atoi(argSlic[0])
		if err != nil {
			panic("The local port provided for the forwarding rule could not be parsed. Please provide a valid port number.")
		}

		remotePort, err := strconv.Atoi(argSlic[3])
		if err != nil {
			panic("The remote port provided for the forwarding rule could not be parsed. Please provide a valid port number.")
		}

		return &ForwardingRule{
			LocalPort:    localPort,
			RemoteClient: argSlic[1],
			RemoteHost:   argSlic[2],
			RemotePort:   remotePort,
			Valid:        false,
		}
	default:
		panic("Invalid forwarding rule. Please provide a rule in the format: localPort:remoteClient[:remoteHost]:remotePort")
	}
}
