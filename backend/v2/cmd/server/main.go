package main

import (
	"log"
	"net/http"

	"chatserver/pkg/websocket"

	"github.com/gorilla/mux"
)

func main() {

	// create websocket hub
	hub := websocket.NewHub()
	go hub.Run()

	r := mux.NewRouter().StrictSlash(true)
	r.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		websocket.WebSocketHandler(hub, w, r)
	})

	log.Fatal(http.ListenAndServe(":8000", r))
}
