package model

import "time"

type MessageType int

const (
	MessageText   MessageType = 0
	MessageLogout MessageType = -1
)

// Message represents the messages exchanged between clients
type Message struct {
	SenderName string    `json:"sendername,omitempty"`
	Body       string    `json:"body"`
	Type       int       `json:"type,omitempty"`
	SentAt     time.Time `json:"sendat,omitempty"`
	ReceivedAt time.Time `json:"receivedat,omitempty"`
}
