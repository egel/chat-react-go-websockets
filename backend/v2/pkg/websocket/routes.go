package websocket

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Handles websocket requests from the peer
func WebSocketHandler(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := &Client{hub: hub, wsconn: conn, send: make(chan []byte, 256)}
	client.hub.register <- client

	// allow collection of memory referenced by the caller by doing all work in
	// new goroutines
	go client.WritePump()
	go client.ReadPump()
}
