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
    case client := <-h.unregister:
      if _, ok := h.clients[client]; ok {
        delete(h.clients, client)
        client.unregistered <- true
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

func (h *Hub) ListClientsOnline() []string {
  var names []string
  for client := range h.clients {
    names = append(names, client.Name())
  }
  return names
}

func (h *Hub) ListSquaresCoords() map[string](*ClientCoords) {
  coords := make(map[string](*ClientCoords), len(h.clients))

  for client := range h.clients {
    coords[client.Name()] = &ClientCoords{
      name: client.Name(),
      XPos: client.SquareX(),
      YPos: client.SquareY(),
    }
  }
  return coords
}
