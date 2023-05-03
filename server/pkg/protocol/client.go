package protocol

import (
  "time"
  "fmt"
  "sync"
  "context"

  ws "github.com/gorilla/websocket"

  "github.com/bd878/live-connections/meta"
  "github.com/bd878/live-connections/server/pkg/rpc"
  "github.com/bd878/live-connections/server/pkg/messages"
)

type Parent interface {
  Broadcaster
  Space
}

type Client struct {
  conn *ws.Conn

  parent Parent

  disk *rpc.Disk

  name string

  square *messages.Coords

  cursor *messages.Coords

  input *Text

  records *Records

  ctx context.Context

  send chan []byte

  quit chan struct{}

  receive chan *messages.RawMessage
}

func NewClient(conn *ws.Conn) *Client {
  conn.SetReadLimit(MaxPayloadSize)
  conn.SetReadDeadline(time.Now().Add(PongWait))
  conn.SetPongHandler(func(string) error { conn.SetReadDeadline(time.Now().Add(PongWait)); return nil })

  square := &messages.Coords{}
  cursor := &messages.Coords{}
  input := &Text{}
  records := &Records{
    List: make([](*messages.Record), 0, 100),
  }

  return &Client{
    conn: conn,
    square: square,
    cursor: cursor,
    input: input,
    records: records,
    send: make(chan []byte, 256),
    quit: make(chan struct{}),
    receive: make(chan *messages.RawMessage),
  }
}

func (c *Client) Send() chan []byte {
  return c.send
}

func (c *Client) Quit() chan struct{} {
  return c.quit
}

func (c *Client) Done() <-chan struct{} {
  return c.ctx.Done()
}

func (c *Client) ParentName() string {
  v, ok := c.parent.(Named)
  if !ok {
    meta.Log().Warn("v is not a named")
    return ""
  }
  return v.Name()
}

func (c *Client) Name() string {
  return c.name
}

func (c *Client) SquareX() float32 {
  return c.square.XPos
}

func (c *Client) SquareY() float32 {
  return c.square.YPos
}

func (c *Client) InputText() string {
  return c.input.Value
}

func (c *Client) Records() [](*messages.Record) {
  return c.records.List
}

func (c *Client) Record() *messages.Record {
  return c.records.Selected
}

func (c *Client) RecordID() int32 {
  return c.Record().CreatedAt
}

func (c *Client) Disk() *rpc.Disk {
  return c.disk
}

func (c *Client) SetDisk(d *rpc.Disk) {
  c.disk = d
}

func (c *Client) SetName(name string) {
  c.name = name
}

func (c *Client) SetParent(p Parent) {
  c.parent = p
}

func (c *Client) SetSelectedRecord(r *messages.Record) {
  c.records.Selected = r
}

func (c *Client) SetRecords(rs [](*messages.Record)) {
  c.records.List = rs
}

func (c *Client) FindRecord(recordID int32) *messages.Record {
  var found *messages.Record
  for i := 0; i < len(c.records.List) && found == nil; i++ {
    if c.records.List[i].ID == recordID {
      found = c.records.List[i]
    }
  }
  return found
}

func (c *Client) SetSquareX(XPos float32) {
  c.square.XPos = XPos
}

func (c *Client) SetSquareY(YPos float32) {
  c.square.YPos = YPos
}

func (c *Client) SetCursorX(XPos float32) {
  c.cursor.XPos = XPos
}

func (c *Client) SetCursorY(YPos float32) {
  c.cursor.YPos = YPos
}

func (c *Client) SetInputText(text string) {
  c.input.Value = text
}

func (c *Client) NewRecord() *messages.Record {
  updatedAt := int32(time.Now().Unix())
  createdAt := updatedAt
  id := createdAt

  rec := &messages.Record{
    Value: "",
    ID: id,
    UpdatedAt: updatedAt,
    CreatedAt: createdAt,
  }

  c.records.List = append(c.records.List, rec)
  return rec
}

func (c *Client) close() {
  c.conn.Close()
  close(c.send)
  close(c.quit)
  close(c.receive)
}

func (c *Client) Run(ctx context.Context) {
  defer c.close()

  innerCtx, cancel := context.WithCancel(ctx)
  c.ctx = innerCtx

  go func() {
    <-c.Quit()
    cancel()
  }()

  var wg sync.WaitGroup
  go func() {
    defer wg.Done()
    wg.Add(1)
    c.receiveLoop()
  }()
  go func() {
    defer wg.Done()
    wg.Add(1)
    c.sendLoop()
  }()

  wg.Wait()
}

func (c *Client) receiveLoop() {
  meta.Log().Debug(c.Name(), "launch reading loop")
  defer meta.Log().Debug(c.Name(), "exit reading loop")

  go func() {
    for {
      select {
      case <-c.Done():
        return

      default:
        _, r, err := c.conn.NextReader()
        if err != nil {
          meta.Log().Warn(c.Name(), "failed to obtain next reader", err)
          continue
        }

        rawMessage := messages.NewRawMessage()
        _, err = rawMessage.ReadFrom(r)
        if err != nil {
          meta.Log().Warn(c.Name(), "failed to read message", err)
          continue
        }

        c.receive <- rawMessage
      }
    }
  }()

  for {
    select {
    case <-c.Done():
      return

    case rm := <- c.receive:
      message, err := rm.Decode()
      if err != nil {
        meta.Log().Warn(c.Name(), "failed to decode message:", err)
        continue
      }

      switch message := message.(type) {
      case *messages.AuthMessage:
        c.Restore()
        c.parent.Join(c)
        defer func() {
          c.Save()
          c.parent.Lose(c)
        }()

      case *messages.TextMessage:
        message.SetUser(c.Name())

        c.SetInputText(message.Str)

        c.parent.Broadcast() <- message.Encode()

      case *messages.MouseMoveMessage:
        message.SetUser(c.Name())

        c.SetCursorX(message.XPos)
        c.SetCursorY(message.YPos)

        c.parent.Broadcast() <- message.Encode()

      case *messages.SquareMoveMessage:
        message.SetUser(c.Name())

        c.SetSquareX(message.XPos)
        c.SetSquareY(message.YPos)

        c.parent.Broadcast() <- message.Encode()

      case *messages.SelectRecordMessage:
        // TODO: implement

      case *messages.AddRecordMessage:
        c.SetSelectedRecord(c.NewRecord())

        responseMessage := messages.NewTitlesListMessage(c.Name(), c.Records())

        c.parent.Broadcast() <- responseMessage.Encode()

      default:
        meta.Log().Warn("unknown event")
      }
    }
  }
}

func (c *Client) sendLoop() {
  meta.Log().Debug(c.Name(), "launch writing loop")

  ticker := time.NewTicker(PingPeriod)

  defer func() {
    meta.Log().Debug(c.Name(), "exit writing loop")
    ticker.Stop()
  }()

  for {
    select {
    case <-c.Done():
      return

    case bytes := <-c.Send():
      c.conn.SetWriteDeadline(time.Now().Add(WriteWait))

      err := c.conn.WriteMessage(ws.BinaryMessage, append(bytes, []byte{'\n'}...))
      if err != nil {
        meta.Log().Warn("failed to write bytes")
        c.Quit() <- struct{}{}
        return
      }

    case <-ticker.C:
      c.conn.SetWriteDeadline(time.Now().Add(WriteWait))

      if err := c.conn.WriteMessage(ws.PingMessage, nil); err != nil {
        meta.Log().Warn(fmt.Sprintf("ping message writing failed, err: %v\n", err))
        c.Quit() <- struct{}{}
        return
      }
    }
  }
}

func (c *Client) Restore() {
  records, err := c.Disk().ListTitles(context.TODO(), c.ParentName(), c.Name())
  if err != nil {
    meta.Log().Error(c.ParentName(), c.Name(), "has no records yet")
    return
  }
  c.SetRecords(records)

  record, err := c.Disk().ReadSelectedTitle(context.TODO(), c.ParentName(), c.Name())
  if err != nil {
    meta.Log().Error("failed to read selected record")
  } else {
    c.SetSelectedRecord(record)
  }

  text, err := c.Disk().ReadText(context.TODO(), c.ParentName(), c.Name(), c.RecordID())
  if err != nil {
    meta.Log().Error("failed to restore client input")
    return
  }
  c.SetInputText(text)

  coords, err := c.Disk().ReadSquareCoords(context.TODO(), c.ParentName(), c.Name())
  if err != nil {
    meta.Log().Error("failed to read client square coords")
    return
  }
  c.SetSquareX(coords.XPos)
  c.SetSquareY(coords.YPos)
}

func (c *Client) Save() {
  c.Disk().WriteSquareCoords(context.TODO(), c.ParentName(), c.Name(), c.SquareX(), c.SquareY())
  c.Disk().WriteText(context.TODO(), c.ParentName(), c.Name(), c.RecordID(), c.InputText())
  c.Disk().SelectTitle(context.TODO(), c.ParentName(), c.Name(), c.RecordID())
}
