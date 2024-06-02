package ws

import (
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
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
		log.Info().
			Any("Origin", r.Header["Origin"]).
			Msg("Origin not accepted")
		return false
	}

	// upgrade this connection to a WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error().Err(err).Msg("WS Connection upgrade failed")
		return
	}
	log.Info().
		Str("Origin", r.Header["Origin"][0]).
		Msg("client connected")

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
