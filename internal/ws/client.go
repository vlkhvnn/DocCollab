package ws

import (
	"log"

	"github.com/gorilla/websocket"
)

// Client represents a single WebSocket connection.
type Client struct {
	Conn *websocket.Conn
	Send chan []byte
}

func (c *Client) ReadPump(room *Room) {
	defer func() {
		room.Unregister <- c
		c.Conn.Close()
	}()
	for {
		_, messageBytes, err := c.Conn.ReadMessage()
		if err != nil {
			log.Printf("Read error: %v", err)
			break
		}

		// Here we assume that clients send "update" messages with the full text.
		room.Broadcast <- BroadcastMessage{
			Sender: c,
			Data:   messageBytes,
		}
	}
}

func (c *Client) WritePump() {
	defer c.Conn.Close()
	for message := range c.Send {
		if err := c.Conn.WriteMessage(websocket.TextMessage, message); err != nil {
			log.Printf("Write error: %v", err)
			break
		}
	}
}
