package api

import (
	"net/http"

	"github.com/CanadianCommander/gopherproxy/cmd/gopherproxyserver/proxy"
	"github.com/CanadianCommander/gopherproxy/internal/logging"
	"github.com/CanadianCommander/gopherproxy/internal/proxcom"
	proxylib "github.com/CanadianCommander/gopherproxy/internal/proxy"
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

		channelName := context.Query(proxylib.ChannelParam)
		if channelName == "" {
			logging.Get().Warn("Incoming connection did not specify a channel")
			context.Status(http.StatusBadRequest)
			return
		}

		clientName := context.Query(proxylib.ClientName)
		if clientName == "" {
			logging.Get().Warn("Incoming connection did not specify a client name")
			context.Status(http.StatusBadRequest)
			return
		}

		client, err := proxylib.UpgradeConnection(context, proxylib.ProxyClientSettings{
			Name:     clientName,
			Channel:  channelName,
			Password: context.GetHeader(proxylib.AuthorizationHeader),
		})

		if err != nil {
			logging.Get().Errorw("Failed to upgrade connection to websocket",
				"error", err)
			context.Status(http.StatusInternalServerError)
		} else {
			err = proxy.Manager.AddEndpoint(client)
			switch err.(type) {
			case *proxylib.AuthenticationError:
				logging.Get().Warnw("Failed to add endpoint to manager. Authentication Error", "error", err.Error())
				client.Write(*proxcom.NewCriticalErrorPacket(err))
			case nil: // no error
			default:
				logging.Get().Errorw("Failed to add endpoint to manager. Unexpected Error ", "error", err.Error())
			}
		}
	} else {
		context.Status(http.StatusBadRequest)
	}
}
