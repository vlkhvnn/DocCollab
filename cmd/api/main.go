package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/vlkhvnn/DocCollab/internal/env"
	"github.com/vlkhvnn/DocCollab/internal/ws"
)

// upgrader is used to upgrade HTTP connections to WebSocket connections.
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func main() {
	cfg := config{
		addr: env.GetString("ADDR", ":8080"),
		hub:  ws.NewHub(),
	}
	app := &application{config: cfg}
	mux := app.mount()
	log.Fatal(app.run(mux))
}
