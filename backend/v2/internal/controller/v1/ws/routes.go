package ws

import (
	"net/http"

	"github.com/gorilla/websocket"
)

const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  socketBufferSize,
	WriteBufferSize: socketBufferSize,
}

// TODO: this is static but can be read from envs
var trustedOriginsList []string = []string{"localhost", "127.0.0.1"}

// WebsocketHandler handles websocket requests from the peer
func WebsocketHandler(hub *Hub, w http.ResponseWriter, r *http.Request) {
	// INFO: temporarily enabling to accept connections from all origins
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	serveWs(hub, w, r)
}
