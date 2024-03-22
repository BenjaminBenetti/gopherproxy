package api

import (
	"net/http"

	"github.com/CanadianCommander/gopherproxy/cmd/gopherproxyserver/proxy"
	"github.com/CanadianCommander/gopherproxy/internal/logging"
	"github.com/CanadianCommander/gopherproxy/internal/websocket"
	"github.com/gin-gonic/gin"
)

const (
	ConnectionListenerRoute = "ws/connect"
)

// ============================================
// Endpoints
// ============================================

// Websocket connect
func ConnectionListen(context *gin.Context) {
	// Upgrade the connection to a websocket
	if context.IsWebsocket() {
		logging.Get().Infow("Incoming websocket connection",
			"remoteAddr", context.Request.RemoteAddr,
			"isWebSocket", context.IsWebsocket(),
		)

		client, err := websocket.UpgradeConnection(context)
		if err != nil {
			logging.Get().Errorw("Failed to upgrade connection to websocket",
				"error", err)
			context.Status(http.StatusInternalServerError)
		} else {
			proxy.Manager.AddEndpoint("whohmmmm", client)
		}
	} else {
		context.Status(http.StatusBadRequest)
	}
}
