// internal/ws/hub.go
package ws

import "sync"

// Hub manages multiple document rooms.
type Hub struct {
	Rooms map[string]*Room
	Mu    sync.Mutex
}

// NewHub creates a new Hub instance.
func NewHub() *Hub {
	return &Hub{
		Rooms: make(map[string]*Room),
	}
}

// GetRoom retrieves an existing room for the given document ID,
// or creates a new one if it doesn't exist.
func (h *Hub) GetRoom(docID string) *Room {
	h.Mu.Lock()
	defer h.Mu.Unlock()
	room, exists := h.Rooms[docID]
	if !exists {
		room = NewRoom(docID)
		h.Rooms[docID] = room
		// Start the room's event loop.
		go room.Run()
	}
	return room
}
