package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gorilla/websocket"
	ws "github.com/vlkhvnn/DocCollab/internal/ws"
)

// NewHub creates and returns a new Hub.
func NewHub() *ws.Hub {
	return &ws.Hub{
		clients:    make(map[*ws.Client]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *ws.Client),
		unregister: make(chan *ws.Client),
	}
}

// Run listens for register, unregister, and broadcast events.
func (h *ws.Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			h.mu.Unlock()
			log.Printf("New client registered. Total clients: %d", len(h.clients))
		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
				log.Printf("Client unregistered. Total clients: %d", len(h.clients))
			}
			h.mu.Unlock()
		case message := <-h.broadcast:
			h.mu.Lock()
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
			h.mu.Unlock()
		}
	}
}

// upgrader is used to upgrade HTTP connections to WebSocket connections.
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// serveWs handles WebSocket requests from clients.
func serveWs(hub *ws.Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection: %v", err)
		return
	}

	client := &ws.Client{conn: conn, send: make(chan []byte, 256)}
	hub.register <- client

	// Start goroutines for reading from and writing to the client.
	go client.readPump(hub)
	go client.writePump()
}

// readPump pumps messages from the WebSocket connection to the hub.
func (c *ws.Client) readPump(hub *ws.Hub) {
	defer func() {
		hub.unregister <- c
		c.conn.Close()
	}()
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			log.Printf("Read error: %v", err)
			break
		}
		hub.broadcast <- message
	}
}c
// writePump pumps messages from the hub to the WebSocket connection.
func (c *ws.Client) writePump() {
	defer c.conn.Close()
	for message := range c.send {
		if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
			log.Printf("Write error: %v", err)
			break
		}
	}
}

func main() {
	hub := NewHub()
	go hub.Run()

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/", homeHandler)
	r.Get("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})

	addr := ":8080"
	log.Printf("Server started on %s", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
