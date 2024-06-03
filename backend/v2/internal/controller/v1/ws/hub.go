package ws

import (
	"github.com/rs/zerolog/log"
)

type Hub struct {
	// Registered clients
	clients map[*Client]bool

	// Inbound messages from the clients
	broadcast chan []byte

	// Register requests from the clients
	register chan *Client

	// Unregister requests from clients
	unregister chan *Client
}

// get new Hub pointer
func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			log.Debug().Any("client", client).Msg("register")
			h.clients[client] = true
			h.send([]byte("A new client connected"), client)
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				log.Debug().Any("client", client).Msg("delete client")
				delete(h.clients, client)
				close(client.send)
				h.send([]byte("client disconnected"), client)
			}
		case message := <-h.broadcast:
			// send to all clients
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					log.Info().Any("client", client).Str("message", string(message)).Msg("Sent message")
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}

// send mesage to clients
func (h *Hub) send(msg []byte, ignore *Client) {
	for conn := range h.clients {
		if conn != ignore {
			h.broadcast <- msg
		}
	}
}
