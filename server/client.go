package server

import (
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
)

type Client struct {
	ID   string
	Hub  *Hub
	Conn *websocket.Conn
	Send chan []byte
}

func NewClient(id string, hub *Hub, conn *websocket.Conn) *Client {
	return &Client{
		ID:   id,
		Hub:  hub,
		Conn: conn,
		Send: make(chan []byte, 256),
	}
}

func (c *Client) ReadPump() {
	defer func() {
		log.Printf("❌ disconnected: %s", c.ID)
		c.Hub.Unregister <- c
		c.Conn.Close()
	}()

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			log.Printf("⚠️ read error from %s: %v", c.ID, err)
			break
		}

		log.Printf("📩 RAW MESSAGE from %s: %s", c.ID, string(message))

		var msg Message
		err = json.Unmarshal(message, &msg)
		if err != nil {
			log.Println("❌ JSON parse error:", err)
			continue
		}

		log.Printf("📩 MESSAGE %s -> %s", msg.From, msg.To)

		// 🔥 КЛЮЧЕВАЯ СТРОКА
		c.Hub.Route <- msg
	}
}

func (c *Client) WritePump() {
	defer c.Conn.Close()

	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			err := c.Conn.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				log.Printf("⚠️ write error to %s: %v", c.ID, err)
				return
			}
		}
	}
}
