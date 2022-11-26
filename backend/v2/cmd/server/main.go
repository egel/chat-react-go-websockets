package main

import (
	"fmt"
	"log"
	"net/http"

	"chatserver/internal/pkg/home"
	"chatserver/internal/pkg/ws"

	"github.com/gorilla/mux"
)

const SERVER_PORT = 8000

func main() {

	// create websocket hub
	hub := ws.NewHub()
	go hub.Run()

	r := mux.NewRouter().StrictSlash(true)
	r.HandleFunc("/", home.HomeHandler)
	r.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		ws.WebsocketHandler(hub, w, r)
	})

	log.Printf("Listening on port: %d", SERVER_PORT)

	var addr = fmt.Sprintf(":%d", SERVER_PORT)
	log.Fatal(http.ListenAndServe(addr, r))
}
