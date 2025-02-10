package forwarddisplay

import (
	"time"

	"github.com/CanadianCommander/gopherproxy/cmd/gopherproxyclient/proxy"
	"github.com/CanadianCommander/gopherproxy/internal/logging"
)

// A logging based "UI" for the forward display. Used by other applications that want to integrate with gopherproxy

type ForwardLoggingUi struct {
	ClientManager *proxy.ClientManager
}

// ==================================================
// Constructor
// ==================================================

// Create a new ForwardLoggingUi
func NewForwardLoggingUi(clientManager *proxy.ClientManager) *ForwardLoggingUi {
	return &ForwardLoggingUi{ClientManager: clientManager}
}

// ==================================================
// Public Methods
// ==================================================

func (ui *ForwardLoggingUi) Build() {
	// nothing to build for a logging ui
}

// start printing out current state info on a loop.
func (ui *ForwardLoggingUi) Start() {

	for {
		<-time.After(1 * time.Second)

		ui.printChannelMemberInfo()
	}

}

// ==================================================
// Private Methods
// ==================================================

func (ui *ForwardLoggingUi) printChannelMemberInfo() {
	logging.Get().Infow("", "tag", "ChannelMembers", "Members", ui.ClientManager.StateManager.ChannelMembers)
}
