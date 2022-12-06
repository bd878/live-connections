package main

type Hub struct {
  clients map[*Client]bool

  register chan *Client

  unregister chan *Client

  broadcast chan []byte
}

var hub = Hub{
  clients: make(map[*Client]bool),
  register: make(chan *Client),
  unregister: make(chan *Client),
  broadcast: make(chan []byte, 256),
}

func GetHub() *Hub {
  return &hub
}

func (h *Hub) run() {
  for {
    select {
    case client := <-h.register:
      h.clients[client] = true
      h.broadcast <- []byte("new client")
    case client := <-h.unregister:
      if _, ok := h.clients[client]; ok {
        delete(h.clients, client)
        close(client.send)
        h.broadcast <- []byte("client left")
      }
    case message := <-h.broadcast:
      for client := range h.clients {
        select {
        case client.send <- message:
        default:
          delete(h.clients, client)
          close(client.send)
        }
      }
    }
  }
}