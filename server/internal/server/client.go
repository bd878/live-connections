package server

import (
  "time"
  "log"

  ws "github.com/gorilla/websocket"
)

const MaxPayloadSize int64 = 2 << 10 // 1024 bytes

const pongWait = 60 * time.Second

const pingPeriod = (pongWait * 9) / 10 // a bit less than pongWait

const writeWait = 10 * time.Second

type Client struct {
  conn *ws.Conn
  hub *Hub
  auth chan bool
  send chan []byte

  area string
  name string
}

func NewClient(conn *ws.Conn, hub *Hub, area, name string) *Client {
  conn.SetReadLimit(MaxPayloadSize)
  conn.SetReadDeadline(time.Now().Add(pongWait))
  conn.SetPongHandler(func(string) error { conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

  return &Client{
    conn: conn,
    hub: hub,
    area: area,
    name: name,
    auth: make(chan bool),
    send: make(chan []byte, 256),
  }
}

func (c *Client) ReadLoop() {
  for {
    _, r, err := c.conn.NextReader()
    if err != nil {
      break
    }

    message := NewMessage()
    if _, err := message.ReadFrom(r); err != nil {
      break
    }

    if err := message.Decode(); err != nil {
      break
    }

    switch message.Type() {
    case authMessageType:
      c.area = message.area
      c.name = message.user

      c.hub.register <- c
      defer c.unregister()
    case mouseMoveMessageType:
      c.hub.broadcast <- message.Encode()
    default:
      log.Println("unknown event =", message.Type())
      break
    }
  }
}

func (c *Client) WriteLoop() {
  ticker := time.NewTicker(pingPeriod)

  defer func() {
    c.conn.Close()
    ticker.Stop()
  }()

  for {
    select {
    case bytes := <-c.send:
      c.write(bytes)
    case <-ticker.C:
      c.ping()
    }
  }
}

func (c *Client) write(p []byte) {
  c.conn.SetWriteDeadline(time.Now().Add(writeWait))

  w, err := c.conn.NextWriter(ws.BinaryMessage)
  defer w.Close()
  if err != nil {
    log.Println("obtaining next writer err =", err)
    return
  }

  if _, err := w.Write(p); err != nil {
    log.Println("failed to write bytes")
    return
  }

  n := len(c.send)
  for i := 0; i < n; i++ {
    if _, err := w.Write(<-c.send); err != nil {
      log.Println("failed to write bytes")
      return
    }
  }
}

func (c *Client) ping() {
  c.conn.SetWriteDeadline(time.Now().Add(writeWait))

  if err := c.conn.WriteMessage(ws.PingMessage, nil); err != nil {
    log.Printf("ping message writing failed, err: %v\n", err)
    return
  }

  return
}

func (c *Client) unregister() {
  c.hub.unregister <- c
  // TODO: write coords to disk
}