package websocket

import (
	"fmt"
	"net/http"

	"github.com/CanadianCommander/gopherproxy/internal/logging"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func UpgradeConnection(context *gin.Context) (*ProxyClient, error) {

	var upgrader = websocket.Upgrader{
		ReadBufferSize:  wsReadBufferSize,
		WriteBufferSize: wsWriteBufferSize,
	}

	var wsCon, err = upgrader.Upgrade(context.Writer, context.Request, nil)
	if err != nil {
		logging.Get().Errorw("Failed to upgrade connection to websocket",
			"error", err)
		return nil, err
	}
	logging.Get().Infow("Connection upgraded to websocket",
		"remoteAddr", context.Request.RemoteAddr)

	return newProxyClient(wsCon), nil
}

// NewOutgoingSocket creates a new outgoing websocket connection to the given url
func NewOutgoingSocket(url string, auth string) (*ProxyClient, error) {
	dialer := websocket.Dialer{
		ReadBufferSize:  wsReadBufferSize,
		WriteBufferSize: wsWriteBufferSize,
	}

	wsCon, _, err := dialer.Dial(url, http.Header{"Authorization": []string{fmt.Sprintf("Basic %s", auth)}})
	if err != nil {
		return nil, err
	}
	return newProxyClient(wsCon), nil
}
