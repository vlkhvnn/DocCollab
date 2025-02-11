package main

import (
	"log"
	"net/http"

	"github.com/vlkhvnn/DocCollab/internal/ws"
)

func (app *application) serveWs(w http.ResponseWriter, r *http.Request) {
	// Extract the document ID from query parameters.
	docID := r.URL.Query().Get("docID")
	if docID == "" {
		http.Error(w, "docID parameter missing", http.StatusBadRequest)
		return
	}

	// Upgrade the HTTP connection to a WebSocket.
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}

	// Create a new client.
	client := &ws.Client{
		Conn: conn,
		Send: make(chan []byte, 256),
	}

	// Get the room for the specified document.
	room := app.config.hub.GetRoom(docID)
	// Register the client with the room.
	room.Register <- client

	// Start client read and write pumps, passing in the room.
	go client.ReadPump(room)
	go client.WritePump()
}
