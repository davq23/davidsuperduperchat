package model

import (
	"crypto/rand"
	"encoding/base64"
	"io"
)

// SessionInfo describes the information stored in a session
type SessionInfo struct {
	id       string
	UserID   int
	Username string
}

// GenerateID generates a new session ID
func (si *SessionInfo) GenerateID(length int64) string {
	b := make([]byte, length)

	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		si.id = ""
	} else {
		si.id = base64.URLEncoding.EncodeToString(b)
	}

	return si.id
}

// GetID returns the session ID
func (si *SessionInfo) GetID() string {
	return si.id
}

// SetID sets the session ID
func (si *SessionInfo) SetID(sid string) {
	si.id = sid
}
