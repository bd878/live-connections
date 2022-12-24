package websocket

import (
  "log"
  "bytes"
  "encoding/binary"
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
      // h.broadcast <- h.DoClientsOnlineMessage()
    case client := <-h.unregister:
      if _, ok := h.clients[client]; ok {
        delete(h.clients, client)
        close(client.send)
        // h.broadcast <- h.DoClientsOnlineMessage()
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

func (h *Hub) ListClientsOnline() []string {
  var names []string
  for client := range h.clients {
    names = append(names, client.name)
  }
  return names
}

func (h *Hub) DoClientsOnlineMessage() []byte {
  var err error

  typeBytes := new(bytes.Buffer)
  if err = binary.Write(typeBytes, enc, listClientsOnlineMessageType); err != nil {
    log.Println("error writing message type =", err)
    return []byte{}
  }

  usersBytes := new(bytes.Buffer)
  var users []string = h.ListClientsOnline()
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