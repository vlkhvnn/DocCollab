package websocket

import (
	"encoding/json"
	"log"
	"sync"
	"time"
)

type Room struct {
	ID         string
	Clients    map[*Client]bool
	Broadcast  chan BroadcastMessage // channel for messages from clients
	Register   chan *Client
	Unregister chan *Client
	Mu         sync.Mutex
	Content    string // the shared document text
}

func NewRoom(id string) *Room {
	return &Room{
		ID:         id,
		Clients:    make(map[*Client]bool),
		Broadcast:  make(chan BroadcastMessage),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Content:    "", // start with empty content
	}
}

// Run starts the room's event loop.
func (r *Room) Run() {
	for {
		select {
		case client := <-r.Register:
			r.Mu.Lock()
			r.Clients[client] = true
			r.Mu.Unlock()
			log.Printf("Client joined room %s. Total clients: %d", r.ID, len(r.Clients))
			// When a client joins, immediately send the current content.
			syncMsg := Message{
				Type:      "sync",
				DocID:     r.ID,
				Position:  0,
				Text:      r.Content,
				UserID:    "server",
				Timestamp: time.Now(),
			}
			data, err := json.Marshal(syncMsg)
			if err == nil {
				client.Send <- data
			}

		case client := <-r.Unregister:
			r.Mu.Lock()
			if _, ok := r.Clients[client]; ok {
				delete(r.Clients, client)
				close(client.Send)
				log.Printf("Client left room %s. Total clients: %d", r.ID, len(r.Clients))
			}
			r.Mu.Unlock()

		case bmsg := <-r.Broadcast:
			// Parse the incoming message.
			var msg Message
			if err := json.Unmarshal(bmsg.Data, &msg); err != nil {
				log.Printf("Error parsing message in room %s: %v", r.ID, err)
				continue
			}

			// We expect messages from clients to have type "update".
			if msg.Type == "update" {
				// Update the shared content.
				r.Mu.Lock()
				r.Content = msg.Text
				r.Mu.Unlock()
				log.Printf("Room %s content updated to: %s", r.ID, r.Content)

				// Create a sync message to broadcast the updated content.
				syncMsg := Message{
					Type:      "sync",
					DocID:     r.ID,
					Position:  0,
					Text:      r.Content,
					UserID:    msg.UserID, // optionally note who made the change
					Timestamp: time.Now(),
				}
				syncData, err := json.Marshal(syncMsg)
				if err != nil {
					log.Printf("Error marshalling sync message: %v", err)
					continue
				}

				// Broadcast the sync message to all clients.
				r.Mu.Lock()
				for client := range r.Clients {
					// Send to all clients, including the sender.
					client.Send <- syncData
				}
				r.Mu.Unlock()
			}
		}
	}
}
