package websocket

import (
  "log"
  "encoding/binary"
  "bytes"
)

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
    case message := <-h.broadcast:
      for client := range h.clients {
        select {
        case client.send <- message:
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
  log.Println("[register] current clients =", clientsOnline)

  h.broadcast <- doClientsOnlineMessage(clientsOnline)
}

func (h *Hub) ListClientsOnline() []string {
  var names []string
  for client := range h.clients {
    names = append(names, client.name)
  }
  return names
}

func doClientsOnlineMessage(users []string) []byte {
  var err error

  typeBytes := new(bytes.Buffer)
  if err = binary.Write(typeBytes, enc, listClientsOnlineMessageType); err != nil {
    log.Println("error writing message type =", err)
    return []byte{}
  }

  usersBytes := new(bytes.Buffer)
  for _, user := range users {
    var size uint16 = uint16(len(user))
    if err = binary.Write(usersBytes, enc, size); err != nil {
      log.Println("error writing user name size =", err)
      return []byte{}
    }

    if err = binary.Write(usersBytes, enc, []byte(user)); err != nil {
      log.Println("error writing user name size =", err)
      return []byte{}
    }
  }

  size := typeBytes.Len() + usersBytes.Len()
  sizeBytes := new(bytes.Buffer)
  if err = binary.Write(sizeBytes, enc, uint16(size)); err != nil {
    log.Println("error writing total size =", err)
    return []byte{}
  }

  return bytes.Join(
    [][]byte{
      sizeBytes.Bytes(),
      typeBytes.Bytes(),
      usersBytes.Bytes(),
    },
    []byte{},
  )
}