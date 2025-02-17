package main

import (
	"log"
	"net/http"

	"github.com/vlkhvnn/DocCollab/internal/websocket"
)

func (app *application) serveWs(w http.ResponseWriter, r *http.Request) {
	docID := r.URL.Query().Get("docID")
	if docID == "" {
		http.Error(w, "docID parameter missing", http.StatusBadRequest)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}

	client := &websocket.Client{
		Conn: conn,
		Send: make(chan []byte, 256),
	}

	room := app.config.hub.GetRoom(docID)
	room.Register <- client

	go client.ReadPump(room)
	go client.WritePump()
}
