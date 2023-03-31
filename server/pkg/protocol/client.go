package protocol

import (
  "time"
  "fmt"
  "context"

  ws "github.com/gorilla/websocket"

  "github.com/teralion/live-connections/meta"
  "github.com/teralion/live-connections/server/pkg/messages"
)

const MaxPayloadSize int64 = 512

const pongWait = 60 * time.Second

const pingPeriod = (pongWait * 9) / 10 // a bit less than pongWait

const writeWait = 10 * time.Second

type Client struct {
  conn *ws.Conn

  areaName string
  myName string

  squareXPos float32
  squareYPos float32

  cursorXPos float32
  cursorYPos float32

  textInput string

  registered chan bool
  unregistered chan bool

  send chan []byte
  area *Area
}

func NewClient(conn *ws.Conn, areaName, myName string) *Client {
  conn.SetReadLimit(MaxPayloadSize)
  conn.SetReadDeadline(time.Now().Add(pongWait))
  conn.SetPongHandler(func(string) error { conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

  return &Client{
    conn: conn,
    areaName: areaName,
    myName: myName,
    send: make(chan []byte, 256),
    registered: make(chan bool),
    unregistered: make(chan bool),
  }
}

func (c *Client) Name() string {
  return c.myName
}

func (c *Client) Area() string {
  return c.areaName
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

func (c *Client) SetArea(area *Area) {
  c.area = area
}

func (c *Client) Close() {
  close(c.send)
  close(c.registered)
  close(c.unregistered)
}

func (c *Client) Run(ctx context.Context) {
  quit := make(chan struct{})

  innerCtx, cancel := context.WithCancel(ctx)

  go c.readLoop(innerCtx, quit)
  go c.lifecycleLoop(innerCtx, quit)
  go c.writeLoop(innerCtx, quit)

  <-quit

  cancel()
  c.Close()
}

func (c *Client) readLoop(ctx context.Context, quit chan struct{}) {
  meta.Log().Debug(c.Name(), "launch reading loop")

  for {
    _, r, err := c.conn.NextReader()
    if err != nil {
      meta.Log().Warn("failed to obtain next reader")
      break
    }

    rawMessage, err := messages.ReadFrom(r)
    if err != nil {
      meta.Log().Warn(c.Name(), "failed to read message", err)
      break
    }

    message, err := rawMessage.Decode()
    if err != nil {
      meta.Log().Warn(c.Name(), "failed to decode message:", err)
      break
    }

    switch message := message.(type) {
    case *messages.AuthMessage:
      meta.Log().Debug(c.Name(), "received auth message")

      c.area.register <- c
      defer func() { c.area.unregister <- c }()
    case *messages.TextMessage:
      meta.Log().Debug(c.Name(), "received text input message")

      message.SetUser(c.Name())

      c.textInput = message.Str

      c.area.broadcast <- message.Encode()
    case *messages.MouseMoveMessage:
      meta.Log().Debug(c.Name(), "received mouse move message")

      message.SetUser(c.Name())

      c.cursorXPos = message.XPos
      c.cursorYPos = message.YPos

      c.area.broadcast <- message.Encode()
    case *messages.SquareMoveMessage:
      meta.Log().Debug(c.Name(), "received square move message")

      message.SetUser(c.Name())

      c.squareXPos = message.XPos
      c.squareYPos = message.YPos

      c.area.broadcast <- message.Encode()
    default:
      meta.Log().Warn("unknown event")
      break
    }
  }

  meta.Log().Debug(c.Name(), "exit reading loop")

  quit <- struct{}{}
}

func (c *Client) writeLoop(ctx context.Context, quit chan struct{}) {
  meta.Log().Debug(c.Name(), "launch writing loop")

  ticker := time.NewTicker(pingPeriod)

  defer func() {
    c.conn.Close()
    ticker.Stop()
  }()

  var exit bool
  for {
    select {
    case bytes := <-c.send:
      c.conn.SetWriteDeadline(time.Now().Add(writeWait))

      w, err := c.conn.NextWriter(ws.BinaryMessage)
      if err != nil {
        meta.Log().Warn("obtaining next writer err =", err)
        exit = true
        break
      }
      defer w.Close()

      meta.Log().Debug("write p =", bytes)

      _, err = w.Write(bytes)
      if err != nil {
        meta.Log().Warn("failed to write bytes")
        exit = true
        break
      }
    case <-ticker.C:
      c.conn.SetWriteDeadline(time.Now().Add(writeWait))

      if err := c.conn.WriteMessage(ws.PingMessage, nil); err != nil {
        meta.Log().Warn(fmt.Sprintf("ping message writing failed, err: %v\n", err))
        exit = true
        break
      }
    }

    if exit {
      break
    }
  }

  meta.Log().Debug(c.Name(), "exit writing loop")

  quit <- struct{}{}
}

func (c *Client) lifecycleLoop(ctx context.Context, quit chan struct{}) {
  meta.Log().Debug(c.Name(), "launch lifecycle loop")

  var closed bool = false
  for {
    select {
    case <-c.registered:
      meta.Log().Debug(c.Name(), "client registered")

      clientsOnlineMessage := messages.NewClientsOnlineMessage(c.area.ListClientsOnline())
      c.area.broadcast <- clientsOnlineMessage.Encode()

      squaresCoords := c.area.ListSquaresCoords()
      for name, coords := range squaresCoords {
        squareInitMessage := messages.NewSquareInitMessage(name, coords.XPos, coords.YPos)
        c.area.broadcast <- squareInitMessage.Encode()
      }

      inputTexts := c.area.ListTextsInputs()
      for name, text := range inputTexts {
        textMessage := messages.NewTextMessage(name, text.Str)
        c.area.broadcast <- textMessage.Encode()
      }
    case <-c.unregistered:
      meta.Log().Debug(c.Name(), "client unregistered")
      close(c.send)
      close(c.registered)
      close(c.unregistered)
      closed = true

      clientsOnlineMessage := messages.NewClientsOnlineMessage(c.area.ListClientsOnline())
      c.area.broadcast <- clientsOnlineMessage.Encode()
    }

    if closed {
      break;
    }
  }

  meta.Log().Debug(c.Name(), "exit lifecycle loop")
}