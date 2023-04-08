package protocol

import (
  "context"

  "github.com/teralion/live-connections/meta"
  "github.com/teralion/live-connections/server/pkg/rpc"
  "github.com/teralion/live-connections/server/pkg/messages"
)

type Area struct {
  disk *rpc.Disk
  clients map[string]*Client
  register chan *Client
  unregister chan *Client
  broadcast chan []byte
}

func NewArea(disk *rpc.Disk) *Area {
  return &Area{
    disk: disk,
    register: make(chan *Client, 1),
    unregister: make(chan *Client, 1),
    broadcast: make(chan []byte, 256),
    clients: make(map[string]*Client, 10),
  }
}

func (a *Area) Close() {
  close(a.register)
  close(a.unregister)
  close(a.broadcast)
}

func (a *Area) Run(ctx context.Context) {
  meta.Log().Debug("area is running")
  defer meta.Log().Debug("area stopped running")

  for {
    select {
    case <-ctx.Done():
      return
    case client := <-a.register:
      a.clients[client.Name()] = client

      a.restoreClient(client)
      go a.onRegister(client)
    case client := <-a.unregister:
      if _, ok := a.clients[client.Name()]; ok {
        a.saveClient(client)

        delete(a.clients, client.Name())

        go a.onUnregister(client)
      }
    case bytes := <-a.broadcast:
      for _, client := range a.clients {
        select {
        case client.send <- bytes:
        default:
          a.unregister <- client
        }
      }
    }
  }
}

func (a *Area) ListClientsOnline() []string {
  var names []string
  for _, client := range a.clients {
    names = append(names, client.Name())
  }
  return names
}

func (a *Area) ListSquaresCoords() map[string](*messages.Coords) {
  coords := make(map[string](*messages.Coords), len(a.clients))

  for _, client := range a.clients {
    coords[client.Name()] = messages.NewCoords(client.SquareX(), client.SquareY())
  }
  return coords
}

func (a *Area) ListTextsInputs() map[string](*messages.Text) {
  texts := make(map[string](*messages.Text), len(a.clients))

  for _, client := range a.clients {
    texts[client.Name()] = messages.NewText(client.TextInput())
  }
  return texts
}

func (a *Area) ListTitlesRecords() map[string]([](*messages.Record)) {
  records := make(map[string]([](*messages.Record)), len(a.clients))

  for _, client := range a.clients {
    records[client.Name()] = client.Records()
  }
  return records
}

func (a *Area) saveClient(c *Client) {
  a.disk.WriteSquareCoords(context.TODO(), c.Area(), c.Name(), c.SquareX(), c.SquareY())
  a.disk.WriteText(context.TODO(), c.Area(), c.Name(), c.TextInput())
}

func (a *Area) restoreClient(c *Client) {
  text, err := a.disk.ReadText(context.TODO(), c.Area(), c.Name())
  if err != nil {
    meta.Log().Error("failed to restore client input")
    return
  }
  c.SetTextInput(text)

  coords, err := a.disk.ReadSquareCoords(context.TODO(), c.Area(), c.Name())
  if err != nil {
    meta.Log().Error("failed to read client square coords")
    return
  }
  c.SetSquareCoords(coords.XPos, coords.YPos)
}

// TODO; send to given client only, broadcast is redundant
func (a *Area) onRegister(c *Client) {
  meta.Log().Debug(c.Name(), "client registered")

  clientsOnlineMessage := messages.NewClientsOnlineMessage(a.ListClientsOnline())
  a.broadcast <- clientsOnlineMessage.Encode()

  squaresCoords := a.ListSquaresCoords()
  for name, coords := range squaresCoords {
    squareInitMessage := messages.NewSquareInitMessage(name, coords.XPos, coords.YPos)
    a.broadcast <- squareInitMessage.Encode()
  }

  titlesRecords := a.ListTitlesRecords()
  for name, records := range titlesRecords {
    titlesListMessage := messages.NewTitlesListMessage(name, records)
    a.broadcast <- titlesListMessage.Encode()
  }

  inputTexts := a.ListTextsInputs()
  for name, text := range inputTexts {
    textMessage := messages.NewTextMessage(name, text.Str)
    a.broadcast <- textMessage.Encode()
  }
}

func (a *Area) onUnregister(c *Client) {
  meta.Log().Debug(c.Name(), "client unregistered")

  clientsOnlineMessage := messages.NewClientsOnlineMessage(a.ListClientsOnline())
  a.broadcast <- clientsOnlineMessage.Encode()
}
