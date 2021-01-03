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
	defer func() {
		c.Hub.Unregister <- c
		c.WS.Close()
	}()

	c.WS.SetReadLimit(1024)
	c.WS.SetReadDeadline(time.Now().Add(time.Duration(time.Minute * 3)))

	c.WS.SetPongHandler(func(string) error {
		c.WS.SetReadDeadline(time.Now().Add(time.Duration(time.Minute)))
		return nil
	})

	for {
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
	defer func() {
		c.Hub.Unregister <- c
		c.WS.Close()
	}()

	ticker := time.NewTicker(time.Duration(1 * time.Minute))

	for {

		select {
		case message, ok := <-c.Send:
			c.WS.SetWriteDeadline(time.Now().Add(time.Duration(time.Minute * 5)))

			if ok {
				message.ReceivedAt = time.Now()
				c.WS.WriteJSON(message)
			} else {
				c.Hub.Unregister <- c
				return
			}
		case <-ticker.C:
			c.Hub.Logger.LogChan <- "updating"

			if err := c.WS.WriteMessage(websocket.PingMessage, nil); err != nil {
				c.Hub.Logger.LogChan <- err.Error()
				return
			}

			c.Hub.Update <- c
		}

	}
}
