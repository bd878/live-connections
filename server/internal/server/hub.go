package server

type Hub struct {
  clients map[*Client]bool

  register chan *Client

  unregister chan *Client

  broadcast chan []byte
}

func NewHub() *Hub {
  return &Hub{
    clients: make(map[*Client]bool),
    register: make(chan *Client),
    unregister: make(chan *Client),
    broadcast: make(chan []byte, 256),
  }
}

func (h *Hub) Run() {
  for {
    select {
    case client := <-h.register:
      h.clients[client] = true

      h.broadcastClientsOnline()
    case client := <-h.unregister:
      if _, ok := h.clients[client]; ok {
        delete(h.clients, client)
        close(client.send)

        h.broadcastClientsOnline()
      }
    case bytes := <-h.broadcast:
      for client := range h.clients {
        select {
        case client.send <- bytes:
        default:
          h.unregister <- client
        }
      }
    }
  }
}

func (h *Hub) broadcastClientsOnline() {
  var clientsOnline []string
  clientsOnline = h.ListClientsOnline()

  h.broadcast <- EncodeClientsOnline(clientsOnline)
}

func (h *Hub) ListClientsOnline() []string {
  var names []string
  for client := range h.clients {
    names = append(names, client.name)
  }
  return names
}
