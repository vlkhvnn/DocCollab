package websocket

import (
	"sync"

	"github.com/vlkhvnn/DocCollab/internal/store"
)

// Hub manages multiple document rooms.
type Hub struct {
	Rooms   map[string]*Room
	Mu      sync.Mutex
	Storage *store.Storage
}

func NewHub(storage *store.Storage) *Hub {
	return &Hub{
		Rooms:   make(map[string]*Room),
		Storage: storage,
	}
}

// GetRoom retrieves or creates a room with the given docID.
func (h *Hub) GetRoom(docID string) *Room {
	h.Mu.Lock()
	defer h.Mu.Unlock()
	room, ok := h.Rooms[docID]
	if !ok {
		room = NewRoom(docID, h.Storage)
		h.Rooms[docID] = room
		go room.Run()
	}
	return room
}
