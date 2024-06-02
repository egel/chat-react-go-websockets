package v1

import (
	"chatserver/internal/controller/v1/ws"
	"net/http"

	"github.com/gorilla/mux"
)

func NewHttpRouter() http.Handler {
	// create websocket hub
	hub := ws.NewHub()
	go hub.Run()

	// configuration
	r := mux.NewRouter().StrictSlash(true)

	// routes
	r.HandleFunc("/", HomeHandler)
	r.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		ws.WebsocketHandler(hub, w, r)
	})

	return r
}
