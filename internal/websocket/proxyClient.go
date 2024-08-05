package websocket

import (
	"time"

	"github.com/CanadianCommander/gopherproxy/internal/logging"
	"github.com/CanadianCommander/gopherproxy/internal/proxcom"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type ProxyClient struct {
	Id            uuid.UUID
	WsCon         *websocket.Conn
	InputChannel  chan proxcom.Packet
	OutputChannel chan proxcom.Packet
	CloseChannel  chan bool
	Closed        bool
	Settings      ProxyClientSettings
}

type ProxyClientSettings struct {
	Name     string
	Channel  string
	Password string
}

// ============================================
// Constructors
// ============================================

// newProxyClient creates a new websocket proxy client
// @param wsCon: the websocket connection
func newProxyClient(wsCon *websocket.Conn, settings ProxyClientSettings) *ProxyClient {
	var client = ProxyClient{
		Id:            uuid.New(),
		WsCon:         wsCon,
		InputChannel:  make(chan proxcom.Packet, proxyChannelBufferSize),
		OutputChannel: make(chan proxcom.Packet, proxyChannelBufferSize),
		CloseChannel:  make(chan bool, 1),
		Closed:        false,
		Settings:      settings,
	}

	wsCon.SetCloseHandler(func(code int, text string) error {
		logging.Get().Infow("Websocket connection closed",
			"code", code,
			"text", text)
		return client.Close()
	})

	go client.messagePump()
	go client.writePump()

	return &client
}

// ============================================
// Public Methods
// ============================================

func (client *ProxyClient) Write(packet proxcom.Packet) {
	client.InputChannel <- packet
}

func (client *ProxyClient) Read() (proxcom.Packet, bool) {
	select {
	case packet, ok := <-client.OutputChannel:
		return packet, ok
	case <-client.CloseChannel:
		return proxcom.Packet{}, false
	}
}

func (client *ProxyClient) Close() error {
	client.Closed = true
	close(client.InputChannel)
	close(client.OutputChannel)
	client.CloseChannel <- true
	close(client.CloseChannel)

	client.WsCon.WriteControl(websocket.CloseMessage, nil, time.Now().Add(1000*time.Millisecond))
	return client.WsCon.Close()
}

// ============================================
// Private Methods
// ============================================

// messagePump reads from the websocket connection and writes to the websocket input channel
func (client *ProxyClient) messagePump() {
	logging.Get().Infow("Starting proxy message pump", "RemoteAddr", client.WsCon.RemoteAddr())

	for {
		if client.Closed {
			break
		}

		msgType, message, err := client.WsCon.ReadMessage()
		if err != nil {
			logging.Get().Warn("Failed to read from websocket, likely close. ",
				"error", err)
			break
		}

		if msgType == websocket.BinaryMessage {
			packet, err := proxcom.NewPacketFromBytes(message)
			if err != nil {
				logging.Get().Warn("Failed to decode incoming packet from remote websocket",
					"error", err,
					"remoteAddr", client.WsCon.RemoteAddr())
			} else {
				client.OutputChannel <- *packet
			}
		}
	}

	logging.Get().Infow("Proxy message pump closed", "RemoteAddr", client.WsCon.RemoteAddr())
}

// writePump reads from the websocket input channel and writes to the websocket connection
func (client *ProxyClient) writePump() {
	logging.Get().Infow("Starting proxy write pump", "RemoteAddr", client.WsCon.RemoteAddr())

	for {
		select {
		case packet := <-client.InputChannel:
			bytes, err := packet.ToBytes()
			if err != nil {
				logging.Get().Warn("Failed to encode packet for sending to remote websocket",
					"error", err,
					"remoteAddr", client.WsCon.RemoteAddr())
			} else {
				err = client.WsCon.WriteMessage(websocket.BinaryMessage, bytes)
				if err != nil {
					logging.Get().Error("Failed to write to remote websocket",
						"error", err,
						"remoteAddr", client.WsCon.RemoteAddr())
				}
			}
		case <-client.CloseChannel:
			logging.Get().Infow("Proxy write pump closed", "RemoteAddr", client.WsCon.RemoteAddr())
			return
		}
	}
}
