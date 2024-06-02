package ws

import (
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// TODO: this is static but can be read from envs
var trustedOriginsList []string = []string{"localhost", "127.0.0.1"}

// WebsocketHandler handles websocket requests from the peer
func WebsocketHandler(hub *Hub, w http.ResponseWriter, r *http.Request) {
	// INFO: enabling to accept connections from trusted origins
	upgrader.CheckOrigin = func(r *http.Request) bool {
		for _, element := range trustedOriginsList {
			if strings.Contains(r.Header["Origin"][0], element) {
				return true
			}
		}
		log.Printf("origin not accepted: %s", r.Header["Origin"])
		return false
	}

	// upgrade this connection to a WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	log.Printf("client connected (from %s)\n", r.Header["Origin"][0])

	websocketClient := &Client{
		hub:    hub,
		wsconn: conn,
		send:   make(chan []byte, 256),
	}
	websocketClient.hub.register <- websocketClient

	// allow collection of memory referenced by the caller by doing all work in
	// new goroutines
	go websocketClient.ReadPump()
	go websocketClient.WritePump()
}
