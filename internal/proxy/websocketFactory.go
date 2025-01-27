package proxy

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/CanadianCommander/gopherproxy/internal/logging"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func UpgradeConnection(context *gin.Context, settings ProxyClientSettings) (*ProxyClient, error) {

	var upgrader = websocket.Upgrader{
		ReadBufferSize:  wsReadBufferSize,
		WriteBufferSize: wsWriteBufferSize,
	}

	var wsCon, err = upgrader.Upgrade(context.Writer, context.Request, nil)
	wsCon.SetReadLimit(wsMaxPacketSize)
	if err != nil {
		logging.Get().Errorw("Failed to upgrade connection to websocket",
			"error", err)
		return nil, err
	}
	logging.Get().Infow("Connection upgraded to websocket",
		"remoteAddr", context.Request.RemoteAddr)

	return newProxyClient(wsCon, settings), nil
}

// NewOutgoingSocket creates a new outgoing websocket connection to the given url
func NewOutgoingSocket(url url.URL, settings ProxyClientSettings) (*ProxyClient, error) {
	dialer := websocket.Dialer{
		ReadBufferSize:  wsReadBufferSize,
		WriteBufferSize: wsWriteBufferSize,
	}

	// set query params
	query := url.Query()
	query.Add(ChannelParam, settings.Channel)
	query.Add(ClientName, settings.Name)
	url.RawQuery = query.Encode()

	wsCon, _, err := dialer.Dial(url.String(), http.Header{AuthorizationHeader: []string{fmt.Sprintf("Basic %s", settings.Password)}})
	if err != nil {
		return nil, err
	}
	return newProxyClient(wsCon, settings), nil
}
