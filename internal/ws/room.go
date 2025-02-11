package ws

import (
	"log"
	"sync"
)

// Room represents a document-specific room where clients collaborate.
type Room struct {
	ID         string
	Clients    map[*Client]bool
	Broadcast  chan []byte
	Register   chan *Client
	Unregister chan *Client
	Mu         sync.Mutex
}

// NewRoom creates a new room for a specific document ID.
func NewRoom(id string) *Room {
	return &Room{
		ID:         id,
		Clients:    make(map[*Client]bool),
		Broadcast:  make(chan []byte),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
	}
}

// Run starts the event loop for the room.
func (r *Room) Run() {
	for {
		select {
		case client := <-r.Register:
			r.Mu.Lock()
			r.Clients[client] = true
			r.Mu.Unlock()
			log.Printf("Client joined room %s. Total clients: %d", r.ID, len(r.Clients))
		case client := <-r.Unregister:
			r.Mu.Lock()
			if _, ok := r.Clients[client]; ok {
				delete(r.Clients, client)
				close(client.Send)
				log.Printf("Client left room %s. Total clients: %d", r.ID, len(r.Clients))
			}
			r.Mu.Unlock()
		case message := <-r.Broadcast:
			r.Mu.Lock()
			for client := range r.Clients {
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					delete(r.Clients, client)
				}
			}
			r.Mu.Unlock()
		}
	}
}
