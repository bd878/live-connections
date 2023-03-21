package server

import (
  "github.com/teralion/live-connections/server/internal/meta"
  "github.com/teralion/live-connections/server/internal/rpc"
)

type Hub struct {
  disk *rpc.Disk
  clients map[*Client]bool
  register chan *Client
  unregister chan *Client
  broadcast chan []byte
}

func NewHub(disk *rpc.Disk) *Hub {
  return &Hub{
    disk: disk,
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

      h.restoreClient(client)

      client.registered <- true
    case client := <-h.unregister:
      if _, ok := h.clients[client]; ok {
        h.saveClient(client)

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

func (h *Hub) ListTextsInputs() map[string](*ClientText) {
  texts := make(map[string](*ClientText), len(h.clients))

  for client := range h.clients {
    texts[client.Name()] = &ClientText{
      name: client.Name(),
      text: client.TextInput(),
    }
  }
  return texts
}

func (h *Hub) saveClient(c *Client) {
  h.disk.WriteSquareCoords(c.Area(), c.Name(), c.SquareX(), c.SquareY())
  h.disk.WriteText(c.Area(), c.Name(), c.TextInput())
}

func (h *Hub) restoreClient(c *Client) {
  text, err := h.disk.ReadText(c.Area(), c.Name())
  if err != nil {
    meta.Log().Error("failed to restore client input")
    return
  }
  c.SetTextInput(text)

  coords, err := h.disk.ReadSquareCoords(c.Area(), c.Name())
  if err != nil {
    meta.Log().Error("failed to read client square coords")
    return
  }
  c.SetSquareCoords(coords.XPos, coords.YPos)
}