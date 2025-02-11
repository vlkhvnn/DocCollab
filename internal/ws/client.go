// internal/ws/client.go
package ws

import (
	"encoding/json"
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

		// Optionally decode JSON (for logging/processing).
		var msg Message
		if err := json.Unmarshal(messageBytes, &msg); err != nil {
			log.Printf("JSON unmarshal error: %v", err)
			continue
		}
		log.Printf("Received in room %s: %+v", room.ID, msg)

		// Broadcast the message to the room.
		room.Broadcast <- messageBytes
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
