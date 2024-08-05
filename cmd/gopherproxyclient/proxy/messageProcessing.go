package proxy

import (
	"fmt"
	"os"

	"github.com/CanadianCommander/gopherproxy/internal/logging"
	"github.com/CanadianCommander/gopherproxy/internal/proxcom"
	"github.com/CanadianCommander/gopherproxy/internal/websocket"
)

// ============================================
// Private Methods
// ============================================

// Error
// CriticalError
// ConfigureEndpoint
// SocketConnect
// SocketDisconnect

func messageProcessingLoop(client *websocket.ProxyClient) {
	for {
		packet, ok := client.Read()
		if !ok {
			logging.Get().Info("Proxy connection closed")
			return
		} else {
			switch packet.Type {
			case proxcom.Data:
				handleData(client, packet)
			case proxcom.Error:
				handleError(client, packet)
			case proxcom.CriticalError:
				handleCriticalError(client, packet)
			case proxcom.ConfigureEndpoint:
				handleConfigureEndpoint(client, packet)
			case proxcom.SocketConnect:
				handleSocketConnect(client, packet)
			case proxcom.SocketDisconnect:
				handleSocketDisconnect(client, packet)
			}
		}
	}
}

func handleData(client *websocket.ProxyClient, packet proxcom.Packet) {
	fmt.Printf("Received data packet from %s: %s\n", packet.Source.Name, string(packet.Data))
}

func handleError(client *websocket.ProxyClient, packet proxcom.Packet) {
	logging.Get().Errorw("Received error packet",
		"error", string(packet.Data))

}

func handleCriticalError(client *websocket.ProxyClient, packet proxcom.Packet) {
	logging.Get().Errorw("Received critical error packet",
		"error", string(packet.Data))
	client.Close()
	os.Exit(1)
}

func handleConfigureEndpoint(client *websocket.ProxyClient, packet proxcom.Packet) {
	logging.Get().Infow("Received configure endpoint packet",
		"endpoint", packet.Target)
}

func handleSocketConnect(client *websocket.ProxyClient, packet proxcom.Packet) {
	logging.Get().Infow("Received socket connect packet",
		"endpoint", packet.Target)
}

func handleSocketDisconnect(client *websocket.ProxyClient, packet proxcom.Packet) {
	logging.Get().Infow("Received socket disconnect packet",
		"endpoint", packet.Target)
}
