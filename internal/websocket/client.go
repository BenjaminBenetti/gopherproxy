package websocket

import (
	"github.com/CanadianCommander/gopherproxy/internal/logging"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type WsClient struct {
	WsCon         *websocket.Conn
	InputChannel  chan []byte
	OutputChannel chan []byte
}

// ============================================
// Constructors
// ============================================

func NewClient(context *gin.Context) (*WsClient, error) {

	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	var wsCon, err = upgrader.Upgrade(context.Writer, context.Request, nil)
	if err != nil {
		logging.Get().Errorw("Failed to upgrade connection to websocket",
			"error", err)
		return nil, err
	}
	logging.Get().Infow("Connection upgraded to websocket",
		"remoteAddr", context.Request.RemoteAddr)

	wsCon.SetCloseHandler(func(code int, text string) error {
		logging.Get().Infow("Websocket connection closed",
			"code", code,
			"text", text)
		return nil
	})

	var client = WsClient{
		WsCon:         wsCon,
		InputChannel:  make(chan []byte),
		OutputChannel: make(chan []byte),
	}

	go doEcho(&client)

	return &client, nil
}

// tmp
func doEcho(client *WsClient) {
	defer func() {
		client.WsCon.Close()
	}()

	for {
		_, message, err := client.WsCon.ReadMessage()
		if err != nil {
			logging.Get().Errorw("Failed to read from websocket",
				"error", err)
			break
		}
		logging.Get().Info("Echo: ", string(message))
		client.WsCon.WriteMessage(websocket.TextMessage, message)
	}
}
