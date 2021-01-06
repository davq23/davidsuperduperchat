package chat

import (
	"context"
	"davidws/model"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// Client represents a client in the network
type Client struct {
	Username  string
	SessionID string
	WS        *websocket.Conn
	Send      chan model.Message
	Hub       *Hub
	lock      sync.Mutex
}

func (c *Client) Read() {
	// Client clean up
	defer func() {
		c.Hub.Unregister <- c
		c.WS.Close()
	}()

	c.WS.SetReadLimit(1024)
	c.WS.SetReadDeadline(time.Now().Add(time.Duration(time.Second * 15)))

	c.WS.SetPongHandler(func(string) error {
		c.WS.SetReadDeadline(time.Now().Add(time.Duration(time.Second * 15)))
		return nil
	})

	for {
		// Grab the next message from the broadcast channel
		var msg model.Message

		// Read in a new message as JSON and map it to a Message object
		err := c.WS.ReadJSON(&msg)

		if err != nil {
			break
		}

		msg.SenderName = c.Username
		msg.SentAt = time.Now()

		switch msg.Type {
		case model.MessageText:
			// Send the newly received message to the broadcast channel
			c.Hub.Broadcast <- msg
		case model.MessageLogout:
			c.Hub.Logger.LogChan <- c.SessionID
			_, err = c.Hub.repo.Delete(context.Background(), c.SessionID)
			return
		}
	}

}

func (c *Client) Write() {
	// Client clean up
	defer func() {
		c.Hub.Unregister <- c
		c.WS.Close()
	}()

	tickerPing := time.NewTicker(time.Duration(1 * time.Second))
	tickerUpdate := time.NewTicker(time.Duration(1 * time.Minute))

	for {
		c.WS.SetWriteDeadline(time.Now().Add(time.Duration(time.Second * 15)))

		select {
		// Received messages
		case message, ok := <-c.Send:

			if ok {
				message.ReceivedAt = time.Now()
				c.WS.WriteJSON(message)
			} else {
				c.Hub.Unregister <- c
				return
			}
		// Check update ticker
		case <-tickerUpdate.C:
			c.Hub.Logger.LogChan <- "updating " + c.Username

			if err := c.WS.WriteMessage(websocket.PingMessage, nil); err != nil {
				c.Hub.Logger.LogChan <- err.Error()
				return
			}

			c.Hub.Update <- c
		// Check keep alive ticker
		case <-tickerPing.C:
			if err := c.WS.WriteMessage(websocket.PingMessage, nil); err != nil {
				c.Hub.Logger.LogChan <- err.Error()
				return
			}
		}

	}
}
