package server

import (
  "github.com/teralion/live-connections/server/internal/meta"
)

type Hub struct {
  area string
  clients map[*Client]bool
  register chan *Client
  unregister chan *Client
  broadcast chan []byte
}

func NewHub(area string) *Hub {
  return &Hub{
    area: area,
    clients: make(map[*Client]bool),
    register: make(chan *Client),
    unregister: make(chan *Client),
    broadcast: make(chan []byte, 256),
  }
}

func (h *Hub) Run() {
  meta.Log().Debug("hub is running")

  for {
    select {
    case client := <-h.register:
      h.clients[client] = true
      client.registered <- true

      h.broadcastClientsOnline()
    case client := <-h.unregister:
      if _, ok := h.clients[client]; ok {
        delete(h.clients, client)
        client.unregistered <- true

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
  meta.Log().Debug("broadcast clients online")

  clientsOnline := h.ListClientsOnline()

  if len(clientsOnline) > 0 {
    encodedClients := EncodeClientsOnline(clientsOnline)
    h.broadcast <- encodedClients
  }
}

func (h *Hub) ListClientsOnline() []string {
  var names []string
  for client := range h.clients {
    names = append(names, client.name)
  }
  return names
}
