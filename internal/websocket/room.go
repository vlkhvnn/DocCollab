// internal/websocket/room.go
package websocket

import (
	"context"
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/vlkhvnn/DocCollab/internal/store"
)

type BroadcastMessage struct {
	Sender *Client
	Data   []byte
}

type Room struct {
	ID         string
	Clients    map[*Client]bool
	Broadcast  chan BroadcastMessage
	Register   chan *Client
	Unregister chan *Client
	Mu         sync.Mutex
	Content    string
	Storage    *store.Storage
}

func NewRoom(docID string, storage *store.Storage) *Room {
	return &Room{
		ID:         docID,
		Clients:    make(map[*Client]bool),
		Broadcast:  make(chan BroadcastMessage),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Content:    "",
		Storage:    storage,
	}
}

func (r *Room) Run() {
	for {
		select {
		case client := <-r.Register:
			r.Mu.Lock()
			r.Clients[client] = true
			r.Mu.Unlock()
			log.Printf("Client joined room %s. Total clients: %d", r.ID, len(r.Clients))
			// When a client joins, send the current content.
			syncMsg := map[string]interface{}{
				"type":      "sync",
				"docID":     r.ID,
				"position":  0,
				"text":      r.Content,
				"userID":    "server",
				"timestamp": time.Now().Format(time.RFC3339),
			}
			data, _ := json.Marshal(syncMsg)
			client.Send <- data

		case client := <-r.Unregister:
			r.Mu.Lock()
			if _, ok := r.Clients[client]; ok {
				delete(r.Clients, client)
				close(client.Send)
				log.Printf("Client left room %s. Total clients: %d", r.ID, len(r.Clients))
			}
			r.Mu.Unlock()

		case bmsg := <-r.Broadcast:
			// Assume the message is an "update" message.
			var msg map[string]interface{}
			if err := json.Unmarshal(bmsg.Data, &msg); err != nil {
				log.Printf("Error parsing message in room %s: %v", r.ID, err)
				continue
			}
			if msg["type"] == "update" {
				newContent, ok := msg["text"].(string)
				if ok {
					// Update in-memory content.
					r.Mu.Lock()
					r.Content = newContent
					r.Mu.Unlock()

					// Persist the update to the database.
					// Here, you can update asynchronously if desired.
					go func() {
						if err := r.Storage.Document.UpdateDocument(context.Background(), r.ID, newContent); err != nil {
							log.Printf("Failed to update document %s: %v", r.ID, err)
						}
					}()

					// Broadcast a sync message with the updated content.
					syncMsg := map[string]interface{}{
						"type":      "sync",
						"docID":     r.ID,
						"position":  0,
						"text":      r.Content,
						"userID":    msg["userID"],
						"timestamp": time.Now().Format(time.RFC3339),
					}
					data, err := json.Marshal(syncMsg)
					if err != nil {
						log.Printf("Error marshalling sync message: %v", err)
						continue
					}
					r.Mu.Lock()
					for client := range r.Clients {
						// Optionally, you might skip sending to the sender if desired.
						client.Send <- data
					}
					r.Mu.Unlock()
				}
			}
		}
	}
}
