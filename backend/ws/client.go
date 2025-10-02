package ws

import (
	"log"
	"github.com/gorilla/websocket"
)

type Client struct {
	Hub *Hub
	Conn *websocket.Conn
	Send chan []byte
}

func (c *Client) ReadPump () {
	defer func() {
		c.Hub.Unregister <- c
		c.Conn.Close()
	}()

	for {
        _, message, err := c.Conn.ReadMessage()
        if err != nil {
            if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
                log.Printf("error: %v", err)
            }
            break
        }
        c.Hub.Broadcast <- message
    }
}

func (c *Client) WritePump() {
	defer func() {
        c.Conn.Close()
    }()
    for {
        select {
        case message, ok := <-c.Send:
            if !ok {
                c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
                return
            }
            c.Conn.WriteMessage(websocket.TextMessage, message)
        }
    }
}