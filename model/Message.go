package model

import "time"

// MessageType type of message
type MessageType int

const (
	// MessageText text message
	MessageText MessageType = 0
	// MessageLogout logout message
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
