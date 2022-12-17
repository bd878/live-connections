package websocket

import (
  "time"
  "log"
  "encoding/binary"
  "net/http"

  ws "github.com/gorilla/websocket"
)

const (
  writeWait = 10 * time.Second

  pongWait = 60 * time.Second

  pingPeriod = (pongWait * 9) / 10 // a bit less than pongWait
)

var enc = binary.LittleEndian

var lenWidth = 2

var newline = []byte{'\n'}

var upgrader = ws.Upgrader{
  HandshakeTimeout: 10 * time.Second,
  ReadBufferSize: 512,
  WriteBufferSize: 512,
}

type Client struct {
  conn *ws.Conn

  hub *Hub

  send chan []byte
}

func NewClient(w http.ResponseWriter, r *http.Request, hub *Hub) {
  conn, err := upgrader.Upgrade(w, r, nil)

  if err != nil {
    log.Println("failed to upgrade connection to WebSocket: ", err)
    return;
  }

  client := &Client{conn: conn, hub: hub, send: make(chan []byte, 256)}
  hub.register <- client

  go client.readLoop()
  go client.writeLoop()
}

func (c *Client) readLoop() {
  defer func() {
    c.hub.unregister <- c
  }()

  c.conn.SetReadLimit(MaxPayloadSize)
  c.conn.SetReadDeadline(time.Now().Add(pongWait))
  c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
  for {
    _, r, err := c.conn.NextReader()
    if err != nil {
      log.Println("NextReader err =", err)
      break
    }

    var size uint16
    if err = binary.Read(r, enc, &size); err != nil {
      log.Println("binary.Read size err =", err)
      break
    }

    log.Println("size =", size)
    message := make([]byte, size)
    if err = binary.Read(r, enc, message); err != nil {
      log.Println("binary.Read message err =", err)
      break
    }
    log.Println("message =", message)

    c.hub.broadcast <- message
  }
}

func (c *Client) writeLoop() {
  ticker := time.NewTicker(pingPeriod)

  defer func() {
    c.conn.Close()
    ticker.Stop()
  }()

  for {
    select {
    case message := <-c.send:
      c.conn.SetWriteDeadline(time.Now().Add(writeWait))

      writer, err := c.conn.NextWriter(ws.BinaryMessage)
      if err != nil {
        log.Println("obtaining next writer err =", err)
        return
      }
      writer.Write(message)

      n := len(c.send)
      for i := 0; i < n; i++ {
        message = <-c.send
        writer.Write(message)
      }

      if err := writer.Close(); err != nil {
        log.Printf("writer close err: %v\n", err)
        return
      }
    case <-ticker.C:
      c.conn.SetWriteDeadline(time.Now().Add(writeWait))
      if err := c.conn.WriteMessage(ws.PingMessage, nil); err != nil {
        log.Printf("ping message writing failed, err: %v\n", err)
        return
      }
    }
  }
}