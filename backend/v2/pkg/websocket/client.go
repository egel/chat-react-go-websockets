package websocket

import (
	"bytes"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// time allowed to write a message to the peer
	writeWait = 10 * time.Second

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
	wsconn *websocket.Conn

	// buffered channel for outbound messages
	send chan []byte
}

// read messages from WebSocket connction to the hub
//
// the application runs ReadPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing allowe
// reads from this goroutine
func (c *Client) ReadPump() {
	defer func() {
		c.hub.unregister <- c
		c.wsconn.Close()
	}()

	c.wsconn.SetReadLimit(maxMessageSize)
	c.wsconn.SetReadDeadline(time.Now().Add(pongWait))
	c.wsconn.SetPongHandler(func(string) error {
		c.wsconn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, message, err := c.wsconn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		c.hub.broadcast <- message
	}
}

// pumps messages from the hub to the websocket connection
func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.wsconn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.wsconn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel
				c.wsconn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.wsconn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// add queued chat messages to the current websocket message
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.wsconn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.wsconn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
