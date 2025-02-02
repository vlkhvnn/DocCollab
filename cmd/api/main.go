package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gorilla/websocket"
)

// upgrader is used to upgrade HTTP connections to WebSocket connections.
var upgrader = websocket.Upgrader{
	// Allow connections from any origin (for development purposes).
	// In production, you should validate the origin.
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// echoHandler upgrades the HTTP connection to WebSocket and echoes back received messages.
func echoHandler(w http.ResponseWriter, r *http.Request) {
	// Upgrade the connection to WebSocket.
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection: %v", err)
		return
	}
	defer conn.Close()

	log.Printf("Client connected: %s", conn.RemoteAddr())

	// Optionally, set a read deadline to avoid hanging connections.
	conn.SetReadDeadline(time.Now().Add(60 * time.Second))

	// Read messages from the client.
	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Read error: %v", err)
			break
		}
		log.Printf("Received: %s", message)

		// Echo the message back to the client.
		if err := conn.WriteMessage(messageType, message); err != nil {
			log.Printf("Write error: %v", err)
			break
		}
	}
}

// homeHandler provides a simple home page.
func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "WebSocket server is running. Connect to /ws")
}

func main() {
	// Create a new chi router.
	r := chi.NewRouter()

	// Use some basic middleware.
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Define HTTP routes.
	r.Get("/", homeHandler)
	r.Get("/ws", echoHandler)

	// Start the HTTP server.
	addr := ":8080"
	log.Printf("Server started on %s", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
