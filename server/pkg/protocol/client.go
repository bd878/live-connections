package protocol

import (
  "time"
  "fmt"
  "context"

  ws "github.com/gorilla/websocket"

  "github.com/bd878/live-connections/meta"
  "github.com/bd878/live-connections/server/pkg/messages"
)

const MaxPayloadSize int64 = 512

const pongWait = 60 * time.Second

const pingPeriod = (pongWait * 9) / 10 // a bit less than pongWait

const writeWait = 10 * time.Second

var newline = []byte{'\n'}

type Client struct {
  conn *ws.Conn

  areaName string
  myName string

  squareXPos float32
  squareYPos float32

  cursorXPos float32
  cursorYPos float32

  textInput string

  records []*messages.Record

  record *messages.Record

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
    records: make([](*messages.Record), 0, 100),
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

func (c *Client) Records() [](*messages.Record) {
  return c.records
}

func (c *Client) Record() *messages.Record {
  return c.record
}

func (c *Client) RecordId() int32 {
  return c.Record().CreatedAt
}

func (c *Client) SetRecord(r *messages.Record) {
  c.record = r
}

func (c *Client) SetRecords(rs [](*messages.Record)) {
  c.records = rs
}

func (c *Client) FindRecord(createdAt int32) *messages.Record {
  var found *messages.Record
  for i := 0; i < len(c.records) && found == nil; i++ {
    if c.records[i].CreatedAt == createdAt {
      found = c.records[i]
    }
  }
  return found
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

func (c *Client) AddNewRecord() *messages.Record {
  updatedAt := int32(time.Now().Unix())
  createdAt := updatedAt

  rec := &messages.Record{
    Value: "",
    UpdatedAt: updatedAt,
    CreatedAt: createdAt,
  }

  c.records = append(c.records, rec)
  return rec
}

func (c *Client) Close() {
  meta.Log().Debug("close client")

  close(c.send)
  close(c.registered)
  close(c.unregistered)
}

func (c *Client) Run(ctx context.Context) {
  defer c.Close()

  quit := make(chan struct{})

  go func() {
    <-ctx.Done()
    quit <- struct{}{}
  }()

  go c.readLoop(quit)
  go c.writeLoop(quit)

  <-quit
}

func (c *Client) readLoop(quit chan struct{}) {
  meta.Log().Debug(c.Name(), "launch reading loop")
  defer meta.Log().Debug(c.Name(), "exit reading loop")

  for {
    _, r, err := c.conn.NextReader()
    if err != nil {
      meta.Log().Warn("failed to obtain next reader")
      quit <- struct{}{}
      return
    }

    rawMessage, err := messages.ReadFrom(r)
    if err != nil {
      meta.Log().Warn(c.Name(), "failed to read message", err)
      continue
    }

    message, err := rawMessage.Decode()
    if err != nil {
      meta.Log().Warn(c.Name(), "failed to decode message:", err)
      continue
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
    case *messages.SelectRecordMessage:
      meta.Log().Debug(c.Name(), "add record message")

      found := c.FindRecord(message.CreatedAt)
      if found != nil {
        meta.Log().Debug("selected record", found.CreatedAt)
      }
    case *messages.AddRecordMessage:
      meta.Log().Debug(c.Name(), "add record message")

      rec := c.AddNewRecord()
      c.SetRecord(rec)
      responseMessage := messages.NewTitlesListMessage(c.Name(), c.Records())

      c.area.broadcast <- responseMessage.Encode()
    default:
      meta.Log().Warn("unknown event")
    }
  }
}

func (c *Client) writeLoop(quit chan struct{}) {
  meta.Log().Debug(c.Name(), "launch writing loop")

  ticker := time.NewTicker(pingPeriod)

  defer func() {
    meta.Log().Debug(c.Name(), "exit writing loop")
    c.conn.Close()
    ticker.Stop()
  }()

  for {
    select {
    case <-quit:
      return
    case bytes := <-c.send:
      c.conn.SetWriteDeadline(time.Now().Add(writeWait))

      meta.Log().Debug("write p =", bytes)

      err := c.conn.WriteMessage(ws.BinaryMessage, append(bytes, newline...))
      if err != nil {
        meta.Log().Warn("failed to write bytes")
        quit <- struct{}{}
        return
      }
    case <-ticker.C:
      c.conn.SetWriteDeadline(time.Now().Add(writeWait))

      if err := c.conn.WriteMessage(ws.PingMessage, nil); err != nil {
        meta.Log().Warn(fmt.Sprintf("ping message writing failed, err: %v\n", err))
        quit <- struct{}{}
        return
      }
    }
  }
}
