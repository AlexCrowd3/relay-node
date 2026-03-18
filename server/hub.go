package server

import (
	"encoding/json"
	"log"
)

type Hub struct {
	Clients    map[string]*Client
	Register   chan *Client
	Unregister chan *Client
	Route      chan Message
}

func NewHub() *Hub {
	return &Hub{
		Clients:    make(map[string]*Client),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Route:      make(chan Message),
	}
}

func (h *Hub) Run() {
	for {
		select {

		case client := <-h.Register:
			h.Clients[client.ID] = client
			log.Printf("✅ registered: %s", client.ID)

		case client := <-h.Unregister:
			if _, ok := h.Clients[client.ID]; ok {
				delete(h.Clients, client.ID)
				close(client.Send)
				log.Printf("❌ unregistered: %s", client.ID)
			}

		case msg := <-h.Route:
			target, ok := h.Clients[msg.To]

			if ok {
				data, err := json.Marshal(msg)
				if err != nil {
					log.Println("❌ marshal error:", err)
					continue
				}

				target.Send <- data
				log.Printf("➡️ routed %s -> %s", msg.From, msg.To)

			} else {
				log.Printf("❌ target not found: %s", msg.To)
			}
		}
	}
}
