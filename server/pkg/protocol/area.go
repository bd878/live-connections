package protocol

import (
  "context"
  "errors"
  "fmt"

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
    clients: make(map[string]*Client, 10),
  }
}

func (a *Area) Close() {
  close(a.register)
  close(a.unregister)
  close(a.broadcast)
}

func (a *Area) Attach(c *Client) {
  a.clients[c.Name()] = c
}

func (a *Area) Run(ctx context.Context) {
  meta.Log().Debug("area is running")

  for {
    select {
    case client := <-a.register:
      a.clients[client.Name()] = client

      a.restoreClient(client)

      client.registered <- true
    case client := <-a.unregister:
      if _, ok := a.clients[client.Name()]; ok {
        a.saveClient(client)

        delete(a.clients, client.Name())
        client.unregistered <- true
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

func (a *Area) GetClients() map[string]*Client {
  return a.clients
}

func (a *Area) GetClient(name string) (*Client, error) {
  if _, has := a.clients[name]; has {
    return a.clients[name], nil
  } else {
    return nil, errors.New(fmt.Sprintf("no client %s", name))
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
    coords[client.Name()] = &messages.Coords{
      XPos: client.SquareX(),
      YPos: client.SquareY(),
    }
  }
  return coords
}

func (a *Area) ListTextsInputs() map[string](*messages.Text) {
  texts := make(map[string](*messages.Text), len(a.clients))

  for _, client := range a.clients {
    texts[client.Name()] = &messages.Text{
      Str: client.TextInput(),
    }
  }
  return texts
}

func (a *Area) saveClient(c *Client) {
  a.disk.WriteSquareCoords(c.Area(), c.Name(), c.SquareX(), c.SquareY())
  a.disk.WriteText(c.Area(), c.Name(), c.TextInput())
}

func (a *Area) restoreClient(c *Client) {
  text, err := a.disk.ReadText(c.Area(), c.Name())
  if err != nil {
    meta.Log().Error("failed to restore client input")
    return
  }
  c.SetTextInput(text)

  coords, err := a.disk.ReadSquareCoords(c.Area(), c.Name())
  if err != nil {
    meta.Log().Error("failed to read client square coords")
    return
  }
  c.SetSquareCoords(coords.XPos, coords.YPos)
}