package chat

import (
	"context"
	"davidws/model"
	"davidws/repo"
	"davidws/utils"
	"time"
)

// Hub is the main component of the chat
type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	Broadcast chan model.Message

	// Register requests from the clients.
	Register chan *Client

	// Unregister requests from clients.
	Unregister chan *Client

	Update chan *Client

	Logger *utils.Logger

	repo repo.SessionRepo
}

// NewHub creates a new Hub
func NewHub(repo repo.SessionRepo, logger *utils.Logger) *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		Broadcast:  make(chan model.Message),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Update:     make(chan *Client),

		Logger: logger,
		repo:   repo,
	}
}

// Run runs the chat app
func (h *Hub) Run(ctx context.Context) {
	for {
		select {
		// Register client
		case client := <-h.Register:
			_, exist := h.clientAuth(ctx, client.SessionID)

			if exist {
				h.clients[client] = true
			} else {
				close(client.Send)
				client.WS.Close()
			}
		// Unregister client
		case client := <-h.Unregister:
			if _, ok := h.clients[client]; ok {
				h.Logger.LogChan <- "Unregistering " + client.Username

				delete(h.clients, client)
				close(client.Send)
				client.WS.Close()
			}
		// Broadcast message to clients
		case message := <-h.Broadcast:
			for client := range h.clients {
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					delete(h.clients, client)
				}
			}
		// Update session back to 15 minute expiration
		case c := <-h.Update:
			err := h.repo.UpdateExpire(context.Background(), c.SessionID, time.Duration(15*time.Minute))

			if err != nil {
				h.Logger.LogChan <- err.Error()
				h.Unregister <- c
			}
		}
	}
}

// clientAuth performs fetches session before registering
func (h *Hub) clientAuth(ctx context.Context, sessionID string) (userID string, ok bool) {
	userIDInterface, err := h.repo.Get(ctx, sessionID)

	if err != nil {
		return "", false
	}

	userID = userIDInterface.(string)

	return userID, true
}
