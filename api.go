package main

import (
	"net/http"

	"github.com/google/uuid"
)

func CursorWebsocketHandler(w http.ResponseWriter, r *http.Request, hub *Hub) {
	conn, err := Upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	client := &Client{
		id:   uuid.New().String(),
		hub:  hub,
		conn: conn,
		send: make(chan []byte, 256),
	}
	client.hub.register <- client

	go client.ReadLoop()
	go client.WriteLoop()
}
