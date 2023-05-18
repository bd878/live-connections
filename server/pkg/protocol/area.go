package protocol

import (
  "context"

  "github.com/bd878/live-connections/meta"
  "github.com/bd878/live-connections/server/pkg/rpc"
  "github.com/bd878/live-connections/server/pkg/messages"
)

type Registry map[string]*Client

type Area struct {
  name string

  disk *rpc.Disk

  registry Registry

  broadcast chan []byte

  running bool
}

func NewArea(name string) *Area {
  return &Area{
    name: name,
    registry: make(Registry, MaxClients),
    broadcast: make(chan []byte, MTU),
  }
}

func (a *Area) Name() string {
  return a.name
}

func (a *Area) Broadcast() chan []byte {
  if !a.Running() {
    meta.Log().Error("area is not running")
    panic("area is in inconsistent state")
  }

  return a.broadcast
}

func (a *Area) Subscribers() *Registry {
  return &a.registry
}

func (a *Area) Disk() *rpc.Disk {
  return a.disk
}

func (a *Area) SetDisk(d *rpc.Disk) {
  a.disk = d
}

func (a *Area) close() {
  close(a.broadcast)
}

func (a *Area) Join(v interface{}) {
  c, ok := v.(*Client)
  if !ok {
    meta.Log().Warn("n is not a Client, cannot join")
    return
  }

  a.registry[c.Name()] = c
  a.onJoin()
}

func (a *Area) Lose(v interface{}) {
  n, ok := v.(Named)
  if !ok {
    meta.Log().Warn("v is not named")
    return
  }

  delete(a.registry, n.Name())
  a.onLeave()
}

func (a *Area) Running() bool {
  return a.running
}

func (a *Area) Run(ctx context.Context) {
  meta.Log().Debug("area is running")
  defer meta.Log().Debug("area stopped running")
  defer func (){ a.running = false }()

  a.running = true

  for {
    select {
    case <-ctx.Done():
      return
    case bytes := <-a.Broadcast():
      for _, c := range a.registry {
        select {
        case c.Send() <- bytes:
        default:
          c.Quit() <- struct{}{}
        }
      }
    }
  }
}

func (a *Area) List() []string {
  var names []string
  for name, _ := range a.registry {
    names = append(names, name)
  }
  return names
}

func (a *Area) listSquaresCoords() map[string](*messages.Coords) {
  coords := make(map[string](*messages.Coords), len(a.registry))

  for name, c := range a.registry {
    coords[name] = messages.NewCoords(c.SquareX(), c.SquareY())
  }

  return coords
}

func (a *Area) listTextsInputs() map[string](*messages.Text) {
  texts := make(map[string](*messages.Text), len(a.registry))

  for name, c := range a.registry {
    texts[name] = messages.NewText(c.Text())
  }
  return texts
}

func (a *Area) listTitlesRecords() map[string]([](*messages.Record)) {
  records := make(map[string]([](*messages.Record)), len(a.registry))

  for name, c := range a.registry {
    records[name] = c.Records()
  }

  return records
}

func (a *Area) onLeave() {
  clientsOnlineMessage := messages.NewClientsOnlineMessage(a.List())
  a.Broadcast() <- clientsOnlineMessage.Encode()
}

func (a *Area) onJoin() {
  clientsOnlineMessage := messages.NewClientsOnlineMessage(a.List())
  a.Broadcast() <- clientsOnlineMessage.Encode()

  squaresCoords := a.listSquaresCoords()
  for name, coords := range squaresCoords {
    squareInitMessage := messages.NewSquareInitMessage(name, coords.XPos, coords.YPos)
    a.Broadcast() <- squareInitMessage.Encode()
  }

  titlesRecords := a.listTitlesRecords()
  for name, records := range titlesRecords {
    titlesListMessage := messages.NewTitlesListMessage(name, records)
    a.Broadcast() <- titlesListMessage.Encode()
  }

  inputTexts := a.listTextsInputs()
  for name, text := range inputTexts {
    textMessage := messages.NewTextMessage(name, text.Value)
    a.Broadcast() <- textMessage.Encode()
  }
}
