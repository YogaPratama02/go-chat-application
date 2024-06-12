package wsdirectmessage

type Hub struct {
	// Registered clients.
	Clients map[*Client]bool

	Broadcast chan Message
	// Register requests from the clients.
	Register chan *Client

	// Unregister requests from clients.
	Unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		Broadcast:  make(chan Message),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Clients:    make(map[*Client]bool),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.Clients[client] = true
		case client := <-h.Unregister:
			if _, ok := h.Clients[client]; ok {
				delete(h.Clients, client)
				close(client.Send)
			}
		case message := <-h.Broadcast:
			for client := range h.Clients {
				if client.SenderId == message.SenderId && client.ReceiverId == message.ReceiverId {
					client.Send <- message
				}

				if client.SenderId == message.ReceiverId && client.ReceiverId == message.SenderId {
					client.Send <- message
				}
			}
		}
	}
}
