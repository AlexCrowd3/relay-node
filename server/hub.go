package server

type Hub struct {
	Clients map[string]*Client

	Register   chan *Client
	Unregister chan *Client

	Route chan Message
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

		case client := <-h.Unregister:

			if _, ok := h.Clients[client.ID]; ok {

				delete(h.Clients, client.ID)

				close(client.Send)
			}

		case msg := <-h.Route:

			target, ok := h.Clients[msg.To]

			if ok {
				target.Send <- msg
			}

		}
	}
}
