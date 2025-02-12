package ws

import "sync"

// Hub manages multiple document rooms.
type Hub struct {
	Rooms map[string]*Room
	Mu    sync.Mutex
}

// NewHub creates a new Hub.
func NewHub() *Hub {
	return &Hub{
		Rooms: make(map[string]*Room),
	}
}

// GetRoom returns an existing room or creates a new one for the given docID.
func (h *Hub) GetRoom(docID string) *Room {
	h.Mu.Lock()
	defer h.Mu.Unlock()
	room, exists := h.Rooms[docID]
	if !exists {
		room = NewRoom(docID)
		h.Rooms[docID] = room
		go room.Run() // Start the room's event loop.
	}
	return room
}
