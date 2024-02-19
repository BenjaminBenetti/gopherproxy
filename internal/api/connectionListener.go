package api

import (
	"net/http"

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
	logging.Get().Infow("Incoming websocket connection",
		"remoteAddr", context.Request.RemoteAddr,
		"isWebSocket", context.IsWebsocket(),
	)

	// Upgrade the connection to a websocket
	if context.IsWebsocket() {
		websocket.NewClient(context)
	} else {
		context.Status(http.StatusBadRequest)
	}
}
