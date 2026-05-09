package main

import (
	"fmt"
	"time"
)

type Hub struct {
	cursorList []Cursor
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
	clients    map[*Client]bool
}

func newHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

func (h *Hub) run() {
	listCursorTicker := time.NewTicker(time.Second)
	for {
		select {
		case client := <-h.register:
			fmt.Println("register client")
			h.clients[client] = true
		case client := <-h.unregister:
			fmt.Println("unregister client")
			delete(h.clients, client)
			close(client.send)
		case message := <-h.broadcast:
			_ = DealMassage(h, message)
		case <-listCursorTicker.C:
			_ = SendCursorCheck(h)
		}
	}
}
