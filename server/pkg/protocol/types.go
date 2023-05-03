package protocol

import (
  "context"

  "github.com/bd878/live-connections/server/pkg/messages"
)

type Text struct {
  Value string
}

type Records struct {
  List []*messages.Record
  Selected *messages.Record
}

type Sender interface {
  Send() chan []byte
  Quit() chan struct{}
}

type Runnable interface {
  Run(ctx context.Context)
}

type Named interface {
  Name() string
}

type Broadcaster interface {
  Broadcast() chan []byte
}

type Squared interface {
  SquareX() float32
  SquareY() float32
  SetSquareX(XPos float32)
  SetSquareY(YPos float32)
}

type Inputable interface {
  InputText() string
  SetInputText(text string)
}

type Catalogable interface {
  Records() []*messages.Record
  SetRecords(rr []*messages.Record)
  FindRecord(recordID int32) *messages.Record
}

type Space interface {
  Join(v Named)
  Lose(v Named)
  List() []string
  Get(n string) (interface{}, error)
}