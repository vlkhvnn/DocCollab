package ws

import "time"

type Message struct {
	Type      string    `json:"type"`
	DocID     string    `json:"docID"`
	Position  int       `json:"position"`
	Text      string    `json:"text"`
	UserID    string    `json:"userID"`
	Timestamp time.Time `json:"timestamp"`
}

type BroadcastMessage struct {
	Sender *Client
	Data   []byte
}
