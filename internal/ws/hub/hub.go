package ws

import (
	"sync"

	client "github.com/vlkhvnn/DocCollab/internal/ws/client"
)

type Hub struct {
	clients    map[*client.Client]bool
	broadcast  chan []byte
	register   chan *client.Client
	unregister chan *client.Client
	mu         sync.Mutex
}
