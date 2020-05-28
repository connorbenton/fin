package socket

// Hub maintains the set of active Clients and Broadcasts messages to the
// Clients.
type Hub struct {
	// Registered Clients.
	Clients map[*Client]bool

	// Inbound messages from the Clients.
	Broadcast chan []byte

	// Reponse reads answers back from frontend
	Response chan []byte

	// Register requests from the Clients.
	Register chan *Client

	// Unregister requests from Clients.
	Unregister chan *Client
}

//ExportHub used in other packages
var ExportHub *Hub

//StartHub is called from main to set up the hub
func StartHub() {
	ExportHub = newHub()
	go ExportHub.run()
}

func newHub() *Hub {
	return &Hub{
		Broadcast:  make(chan []byte),
		Response:   make(chan []byte),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Clients:    make(map[*Client]bool),
	}
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.Register:
			h.Clients[client] = true
		case client := <-h.Unregister:
			if _, ok := h.Clients[client]; ok {
				delete(h.Clients, client)
				close(client.send)
			}
		case message := <-h.Broadcast:
			for client := range h.Clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.Clients, client)
				}
			}
		}
	}
}
