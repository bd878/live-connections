package protocol

import (
  "time"
  "fmt"
  "sync"
  "context"

  "github.com/bd878/live-connections/meta"
  "github.com/bd878/live-connections/server/pkg/rpc"
  "github.com/bd878/live-connections/server/pkg/messages"
)

type Parent interface {
  Broadcaster
  Space
  Named
}

type Disk interface {
  rpc.IdentityDisk
  rpc.SquareDisk
  rpc.TextDisk
  rpc.CatalogDisk
}

type Client struct {
  conn Conn

  parent Parent

  disk Disk

  name string

  square *messages.Coords

  cursor *messages.Coords

  records *Records

  ctx context.Context

  send chan []byte

  quit chan struct{}

  receive chan *messages.RawMessage

  running bool
}

func NewClient(conn Conn) *Client {
  square := &messages.Coords{}
  cursor := &messages.Coords{}
  records := &Records{
    List: make([](*messages.TextRecord), 0, 100),
    Selected: nil,
  }

  return &Client{
    conn: conn,
    square: square,
    cursor: cursor,
    records: records,
    send: make(chan []byte, 256),
    quit: make(chan struct{}),
    receive: make(chan *messages.RawMessage),
  }
}

func (c *Client) Send() chan []byte {
  if !c.Running() {
    meta.Log().Error("client is not running")
    panic("client is in inconsistent state")
  }

  return c.send
}

func (c *Client) Quit() chan struct{} {
  return c.quit
}

func (c *Client) Done() <-chan struct{} {
  return c.ctx.Done()
}

func (c *Client) ParentName() string {
  return c.parent.Name()
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

func (c *Client) Text() string {
  if c.SelectedRecord() == nil {
    meta.Log().Warn("record is not selected")
    return ""
  }

  return c.SelectedRecord().Text.Value
}

func (c *Client) Records() [](*messages.TextRecord) {
  return c.records.List
}

func (c *Client) SelectedRecord() *messages.TextRecord {
  return c.records.Selected
}

func (c *Client) RecordID() int32 {
  return c.SelectedRecord().ID
}

func (c *Client) Disk() Disk {
  return c.disk
}

func (c *Client) SetDisk(d Disk) {
  c.disk = d
}

func (c *Client) SetName(name string) {
  c.name = name
}

func (c *Client) SetParent(v interface{}) {
  p, ok := v.(Parent)
  if !ok {
    meta.Log().Warn("not a parent")
    return
  }
  c.parent = p
}

func (c *Client) Parent() Parent {
  return c.parent
}

func (c *Client) SelectRecord(r *messages.TextRecord) {
  c.records.Selected = r
}

func (c *Client) SetRecords(rs [](*messages.TextRecord)) {
  c.records.List = rs
}

func (c *Client) FindRecord(recordID int32) *messages.TextRecord {
  var found *messages.TextRecord
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

func (c *Client) SetText(text string) {
  if c.SelectedRecord() != nil {
    c.Disk().WriteText(context.TODO(), c.ParentName(), c.Name(), c.SelectedRecord().ID, text)
    // TODO: get selected record
    c.SelectedRecord().Text.Value = text
  } else {
    meta.Log().Warn("record is not selected")
  }
}

func (c *Client) AddNewRecord() *messages.TextRecord {
  rec, err := c.Disk().AddTextRecord(context.TODO(), c.ParentName(), c.Name())
  if err != nil {
    meta.Log().Debug("failed to add text record")
    return nil
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

func (c *Client) Running() bool {
  return c.running
}

func (c *Client) Run(ctx context.Context) {
  defer c.close()
  defer func (){ c.running = false }()

  c.running = true

  innerCtx, cancel := context.WithCancel(ctx)
  c.ctx = innerCtx

  go func() {
    <-c.Quit()
    cancel()
  }()

  var wg sync.WaitGroup
  wg.Add(2)
  go func() {
    defer wg.Done()
    c.receiveLoop()
  }()
  go func() {
    defer wg.Done()
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
          c.Quit() <- struct{}{}
          return
        }

        rawMessage := messages.NewRawMessage()
        _, err = rawMessage.ReadFrom(r)
        if err != nil {
          meta.Log().Warn(c.Name(), "failed to read message", err)
          c.Quit() <- struct{}{}
          return
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

        c.SetText(message.Value)

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
        c.SelectRecord(c.AddNewRecord())

        responseMessage := messages.NewRecordsListMessage(c.Name(), c.Records())

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
      err := c.conn.WriteMessage(BinaryMessage, append(bytes, []byte{'\n'}...))
      if err != nil {
        meta.Log().Warn("failed to write bytes")
        c.Quit() <- struct{}{}
        return
      }

    case <-ticker.C:
      if err := c.conn.WriteMessage(PingMessage, nil); err != nil {
        meta.Log().Warn(fmt.Sprintf("ping message writing failed, err: %v\n", err))
        c.Quit() <- struct{}{}
        return
      }
    }
  }
}

func (c *Client) Restore() {
  records, err := c.Disk().ListTextRecords(context.TODO(), c.ParentName(), c.Name())
  if err != nil {
    meta.Log().Error(c.ParentName(), c.Name(), "has no records yet")
    return
  }
  c.SetRecords(records)

  record, err := c.Disk().GetSelectedRecord(context.TODO(), c.ParentName(), c.Name())
  if err != nil {
    meta.Log().Error("failed to read selected record")
  } else {
    c.SelectRecord(record)
  }

  c.SetText(record.Text.Value)

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
  c.Disk().SelectTextRecord(context.TODO(), c.ParentName(), c.Name(), c.RecordID())
  c.Disk().WriteText(context.TODO(), c.ParentName(), c.Name(), c.RecordID(), c.Text())
}
