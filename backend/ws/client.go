package ws

import (
	"encoding/json"
	"live-code/backend/docker"
	"log"

	"github.com/gorilla/websocket"
)

type Client struct {
	Hub *Hub
	Conn *websocket.Conn
	Send chan []byte
    Manager *docker.Manager
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

        var msg Message
        if err:=json.Unmarshal(message, &msg); err!=nil {
            log.Printf("error unmarshalling message: %v", err)
			continue
        }

        switch msg.Type {
        case "code_update":
            c.Hub.Broadcast <- message
        case "run_code":        
            output, err := c.Manager.RunCode(msg.Payload)
            if err != nil {
				log.Printf("error running code: %v", err)
				output = "Error executing code."
			}

			outputMsg := Message{
				Type:    "run_output",
				Payload: output,
			}
            outputJson, err := json.Marshal(outputMsg)
			if err != nil {
				log.Printf("error marshalling output: %v", err)
				continue
			}
			c.Hub.Broadcast <- outputJson

        }
        
        
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