package server

import (
	"log"

	"github.com/gorilla/websocket"
)

type Client struct {
	ID string

	Hub *Hub

	Conn *websocket.Conn

	Send chan Message
}

func NewClient(id string, hub *Hub, conn *websocket.Conn) *Client {

	return &Client{

		ID: id,

		Hub: hub,

		Conn: conn,

		Send: make(chan Message, 256),
	}
}

func (c *Client) ReadPump() {

	defer func() {

		c.Hub.Unregister <- c

		c.Conn.Close()
	}()

	for {

		var msg Message

		err := c.Conn.ReadJSON(&msg)

		if err != nil {

			log.Println(err)

			break
		}

		c.Hub.Route <- msg
	}
}

func (c *Client) WritePump() {

	defer c.Conn.Close()

	for {

		msg, ok := <-c.Send

		if !ok {

			c.Conn.WriteMessage(websocket.CloseMessage, []byte{})

			return
		}

		err := c.Conn.WriteJSON(msg)

		if err != nil {

			return
		}
	}
}
