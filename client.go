package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

var Upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Client struct {
	id   string
	hub  *Hub
	conn *websocket.Conn
	send chan []byte
}

func (c *Client) ReadLoop() {
	defer func() {
		c.hub.unregister <- c
		_ = c.conn.Close()
	}()
	for {
		_, message, err := c.conn.ReadMessage()
		fmt.Println(fmt.Sprintf("[%s]: %s", c.id, message))
		if err != nil {
			break
		}
		c.hub.broadcast <- message
	}
}

func (c *Client) WriteLoop() {
	defer func() {
		_ = c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				_ = c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
				return
			}
		}

	}
}
