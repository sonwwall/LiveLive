package ws

import (
	"log"
	"net/http"
)

func InitWebsocket() {
	hub := NewHub()

	http.HandleFunc("/ws", NewHandler(hub))

	log.Println("WebSocket server started on :8060")
	log.Fatal(http.ListenAndServe(":8060", nil))
}
