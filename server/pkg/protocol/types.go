package protocol

import (
  "context"
  "io"

  "github.com/bd878/live-connections/server/pkg/messages"
)

type Text struct {
  Value string
}

type Records struct {
  List []*messages.TextRecord
  Selected *messages.TextRecord
}

type Conn interface {
  NextReader() (int, io.Reader, error)
  WriteMessage(messageType int, data []byte) error
  Close() error
}

type Sender interface {
  Send() chan []byte
  Quit() chan struct{}
}

type Runnable interface {
  Running() bool
  Run(ctx context.Context)
}

type Named interface {
  Name() string
}

type Broadcaster interface {
  Broadcast() chan []byte
}

type Children interface {
  SetParent(v interface{})
  Parent() interface{}
}

type Squared interface {
  SquareX() float32
  SquareY() float32
  SetSquareX(XPos float32)
  SetSquareY(YPos float32)
}

type Catalogable interface {
  Records() []*messages.TextRecord
  SelectedRecord() *messages.TextRecord
  SetRecords(rr []*messages.TextRecord)
  FindRecord(recordID int32) *messages.TextRecord
  SelectRecord(recordID int32)
  Text() string
  SetText(txt string)
}

type Space interface {
  Join(v interface{})
  Lose(v interface{})
  List() []string
}
