package main

import (
	"log"
	"net/http"

	"github.com/vlkhvnn/DocCollab/internal/ws"
)

func (app *application) serveWs(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection: %v", err)
		return
	}

	client := &ws.Client{Conn: conn, Send: make(chan []byte, 256)}
	app.config.hub.Register <- client

	// Start goroutines for reading from and writing to the client.
	go client.ReadPump(app.config.hub)
	go client.WritePump()
}
