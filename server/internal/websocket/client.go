package websocket

import (
  "time"
  "log"
  "net/http"
  "bytes"

  ws "github.com/gorilla/websocket"
)

const (
  writeWait = 10 * time.Second

  pongWait = 60 * time.Second

  maxRead = 512 // bytes

  pingPeriod = (pongWait * 9) / 10 // a bit less than pongWait
)

var newline = []byte{'\n'}

var upgrader = ws.Upgrader{
  HandshakeTimeout: 10 * time.Second,
  ReadBufferSize: 1024,
  WriteBufferSize: 1024,
}

type Client struct {
  conn *ws.Conn

  hub *Hub

  send chan []byte
}

func NewClient(w http.ResponseWriter, r *http.Request) {
  conn, err := upgrader.Upgrade(w, r, nil)

  if err != nil {
    log.Println("failed to upgrade connection to WebSocket: ", err)
    return;
  }

  hub := GetHub()
  if err != nil {
    log.Printf("error obtaining hub: %v\n", err)
    return
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

  c.conn.SetReadLimit(maxRead)
  c.conn.SetReadDeadline(time.Now().Add(pongWait))
  c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
  for {
    _, message, err := c.conn.ReadMessage()
    if err != nil {
      if ws.IsUnexpectedCloseError(err, ws.CloseGoingAway) {
        log.Println("err: %v\n", err)
      }
      break;
    }
    message = bytes.TrimSpace(message)
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