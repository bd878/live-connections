package server

import (
  "time"

  ws "github.com/gorilla/websocket"
  "github.com/teralion/live-connections/meta"
  "github.com/teralion/live-connections/server/protocol"
)

const MaxPayloadSize int64 = 512

const pongWait = 60 * time.Second

const pingPeriod = (pongWait * 9) / 10 // a bit less than pongWait

const writeWait = 10 * time.Second

type Client struct {
  conn *ws.Conn
  hub *Hub

  send chan []byte

  registered chan bool
  unregistered chan bool

  area string
  name string

  squareXPos float32
  squareYPos float32

  cursorXPos float32
  cursorYPos float32

  textInput string
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
    registered: make(chan bool),
    unregistered: make(chan bool),
    send: make(chan []byte, 256),
  }
}

func (c *Client) Name() string {
  return c.name
}

func (c *Client) Area() string {
  return c.area
}

func (c *Client) SquareX() float32 {
  return c.squareXPos
}

func (c *Client) SquareY() float32 {
  return c.squareYPos
}

func (c *Client) TextInput() string {
  return c.textInput
}

func (c *Client) SetSquareCoords(XPos, YPos float32) {
  c.squareXPos = XPos
  c.squareYPos = YPos
}

func (c *Client) SetTextInput(text string) {
  c.textInput = text
}

func (c *Client) ReadLoop() {
  meta.Log().Debug(c.name, "launch reading loop")

  for {
    _, r, err := c.conn.NextReader()
    if err != nil {
      meta.Log().Warn("failed to obtain next reader")
      break
    }

    rawMessage, err := protocol.ReadFrom(r)
    if err != nil {
      meta.Log().Warn(c.name, "failed to read message", err)
      break
    }

    message, err := rawMessage.Decode()
    if err != nil {
      meta.Log().Warn(c.name, "failed to decode message:", err)
      break
    }

    switch message := message.(type) {
    case *protocol.AuthMessage:
      meta.Log().Debug(c.name, "received auth message")

      c.hub.register <- c
      defer c.unregister()
    case *protocol.TextMessage:
      meta.Log().Debug(c.name, "received text input message")

      message.SetUser(c.name)

      c.textInput = message.Str

      c.hub.broadcast <- message.Encode()
    case *protocol.MouseMoveMessage:
      meta.Log().Debug(c.name, "received mouse move message")

      message.SetUser(c.name)

      c.cursorXPos = message.XPos
      c.cursorYPos = message.YPos

      c.hub.broadcast <- message.Encode()
    case *protocol.SquareMoveMessage:
      meta.Log().Debug(c.name, "received square move message")

      message.SetUser(c.name)

      c.squareXPos = message.XPos
      c.squareYPos = message.YPos

      c.hub.broadcast <- message.Encode()
    default:
      meta.Log().Warn("unknown event")
      break
    }
  }

  meta.Log().Debug(c.name, "exit reading loop")
}

func (c *Client) LifecycleLoop() {
  meta.Log().Debug(c.name, "launch lifecycle loop")

  var closed bool = false
  for {
    select {
    case <-c.registered:
      meta.Log().Debug(c.name, "client registered")

      clientsOnlineMessage := protocol.NewClientsOnlineMessage(c.hub.ListClientsOnline())
      c.hub.broadcast <- clientsOnlineMessage.Encode()

      squaresCoords := c.hub.ListSquaresCoords()
      for _, coords := range squaresCoords {
        squareInitMessage := protocol.NewSquareInitMessage(coords.name, coords.XPos, coords.YPos)
        c.hub.broadcast <- squareInitMessage.Encode()
      }

      inputTexts := c.hub.ListTextsInputs()
      for _, text := range inputTexts {
        textMessage := protocol.NewTextMessage(text.name, text.text)
        c.hub.broadcast <- textMessage.Encode()
      }
    case <-c.unregistered:
      meta.Log().Debug(c.name, "client unregistered")
      close(c.send)
      close(c.registered)
      close(c.unregistered)
      closed = true

      clientsOnlineMessage := protocol.NewClientsOnlineMessage(c.hub.ListClientsOnline())
      c.hub.broadcast <- clientsOnlineMessage.Encode()
    }

    if closed {
      break;
    }
  }

  meta.Log().Debug(c.name, "exit lifecycle loop")
}

func (c *Client) WriteLoop() {
  meta.Log().Debug(c.name, "launch writing loop")

  ticker := time.NewTicker(pingPeriod)

  defer func() {
    c.conn.Close()
    ticker.Stop()
  }()

  var err error
  for {
    select {
    case bytes := <-c.send:
      err = c.write(bytes)
    case <-ticker.C:
      err = c.ping()
    }

    if err != nil {
      break
    }
  }

  meta.Log().Debug(c.name, "exit writing loop")
}

func (c *Client) write(p []byte) error {
  c.conn.SetWriteDeadline(time.Now().Add(writeWait))

  w, err := c.conn.NextWriter(ws.BinaryMessage)
  if err != nil {
    meta.Log().Warn("obtaining next writer err =", err)
    return err
  }
  defer w.Close()

  meta.Log().Debug("write p =", p)

  if _, err := w.Write(p); err != nil {
    meta.Log().Warn("failed to write bytes")
    return err
  }

  return nil
}

func (c *Client) ping() error {
  c.conn.SetWriteDeadline(time.Now().Add(writeWait))

  if err := c.conn.WriteMessage(ws.PingMessage, nil); err != nil {
    meta.Log().Warn("ping message writing failed, err: %v\n", err)
    return err
  }

  return nil
}

func (c *Client) unregister() {
  meta.Log().Debug(c.name, "unregister")

  c.hub.unregister <- c
  // TODO: write coords to disk
}