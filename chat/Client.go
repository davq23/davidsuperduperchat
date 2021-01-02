package chat

import (
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

	for {
		c.WS.SetReadDeadline(time.Now().Add(time.Duration(time.Minute * 5)))

		// Grab the next message from the broadcast channel
		var msg model.Message

		// Read in a new message as JSON and map it to a Message object
		err := c.WS.ReadJSON(&msg)

		if err != nil {
			c.lock.Lock()
			c.Hub.Unregister <- c
			c.lock.Unlock()
			break
		}

		msg.SenderName = c.Username
		msg.SentAt = time.Now()

		// Send the newly received message to the broadcast channel
		c.Hub.Broadcast <- msg
	}

}

func (c *Client) Write() {
	ticker := time.NewTicker(time.Duration(1 * time.Minute))

	for {
		c.WS.SetWriteDeadline(time.Now().Add(time.Duration(time.Minute * 5)))

		select {
		case message, ok := <-c.Send:
			if ok {
				message.ReceivedAt = time.Now()
				c.WS.WriteJSON(message)
			} else {
				c.Hub.Unregister <- c
			}
		case <-ticker.C:
			c.Hub.Logger.LogChan <- "updating"
			c.lock.Lock()

			if err := c.WS.WriteMessage(websocket.PingMessage, nil); err != nil {
				c.Hub.Unregister <- c
				c.lock.Unlock()
				c.Hub.Logger.LogChan <- err.Error()
				return
			}

			c.Hub.Update <- c
			c.lock.Unlock()
		}

	}
}
