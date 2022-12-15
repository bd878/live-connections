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
    p, r, err := c.conn.NextReader()
    if err != nil {
      log.Printf("NextReader err: %v\n", err)
      break
    }

    var messageType int8
    if err = binary.Read(r, enc, &messageType); err != nil {
      log.Printf("reader.Read err: %v\n", err)
      break
    }

    coords := make([]float32, 2)
    if err := binary.Read(r, enc, &coords); err != nil {
      log.Printf("binary.Read failed: %v\n", err)
      break
    }

    log.Printf("%v, %v : %v\n", p, messageType, coords)

    if err != nil {
      if ws.IsUnexpectedCloseError(err, ws.CloseGoingAway) {
        log.Printf("err: %v\n", err)
      }
      break;
    }
    c.hub.broadcast <- make([]byte, 0)
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
    case message, ok := <-c.send:
      c.conn.SetWriteDeadline(time.Now().Add(writeWait))
      if !ok {
        c.conn.WriteMessage(ws.CloseMessage, []byte{})
        return
      }

      writer, err := c.conn.NextWriter(ws.TextMessage)
      if err != nil {
        log.Printf("obtaining next writer err: %v\n", err)
        return
      }
      writer.Write(message)

      n := len(c.send)
      for i := 0; i < n; i++ {
        writer.Write(newline)
        writer.Write(<-c.send)
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