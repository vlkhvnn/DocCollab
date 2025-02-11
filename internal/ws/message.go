package ws

import "time"

type Message struct {
	Type      string    `json:"type"`
	DocID     string    `json:"docID"`
	Position  int       `json:"position"`
	Text      string    `json:"text"`
	User      string    `json:"user"`
	Timestamp time.Time `json:"timestamp"`
}
