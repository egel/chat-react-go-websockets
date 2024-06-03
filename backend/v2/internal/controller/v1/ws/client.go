package ws

import (
	"bytes"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
)

const (
	// time allowed to write a message to the peer
	writeWait = 600 * time.Second

	// time allowed to read the next pong message from peer
	pongWait = 60 * time.Second

	// send pings to peer with this period. must be less than pongWait
	pingPeriod = (pongWait * 9) / 10

	// maximum messge size allowe from peer (bytes)
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}

	space = []byte{' '}
)

type Client struct {
	// hold Hub that this websocket client belongs
	hub *Hub

	// hold WebSocket Connection
	conn *websocket.Conn

	// channel for outbound messages
	send chan []byte
}

// read messages from WebSocket connction to the hub
//
// the application runs ReadPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing allowe
// reads from this goroutine
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	// c.conn.SetReadLimit(maxMessageSize)
	// c.conn.SetReadDeadline(time.Now().Add(pongWait))
	// c.conn.SetPongHandler(func(string) error {
	// c.conn.SetReadDeadline(time.Now().Add(pongWait))
	// return nil
	// })

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Error().Err(err).Msg("ERROR")
			}
			break
		}
		message = bytes.TrimSpace(bytes.ReplaceAll(message, newline, space))
		c.hub.broadcast <- message
	}
}

// pumps messages from the hub to the websocket connection
func (c *Client) writePump() {
	defer func() {
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			// c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				_ = c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			_ = c.conn.WriteMessage(websocket.TextMessage, message)
		}
	}
}

func serveWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error().Err(err).Msg("WS connection upgrade error")
		http.NotFound(w, r)
		return
	}
	client := &Client{
		hub:  hub,
		conn: conn,
		send: make(chan []byte, maxMessageSize),
	}
	client.hub.register <- client

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump()
	go client.readPump()
}
