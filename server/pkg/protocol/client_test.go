package protocol

import (
  "testing"
  "reflect"
  "io"
  "time"
  "context"
)

type NullConn struct {
  messages chan []byte
}

func NewNullConn() *NullConn {
  return &NullConn{
    messages: make(chan []byte),
  }
}

func (c *NullConn) Messages() chan []byte {
  return c.messages
}

func (c *NullConn) NextReader() (int, io.Reader, error) {
  ch := make(chan struct{})
  <-ch

  return 0, nil, nil
}

func (c *NullConn) WriteMessage(_ int, data []byte) error {
  c.Messages() <- data
  return nil
}

func (c *NullConn) Close() error {
  return nil
}

type NullParent struct {
  name string
  broadcast chan []byte
  clients map[string]Named
}

func NewNullParent() *NullParent {
  return &NullParent{
    name: "test",
    broadcast: make(chan []byte, 256),
    clients: make(map[string]Named, 1),
  }
}

func (p *NullParent) Name() string {
  return p.name
}

func (p *NullParent) Broadcast() chan []byte {
  return p.broadcast
}

func (p *NullParent) Join(v interface{}) {
  n, ok := v.(Named)
  if !ok {
    panic("not named")
  }

  p.clients[n.Name()] = n
}

func (p *NullParent) Lose(v interface{}) {
  n, ok := v.(Named)
  if !ok {
    panic("not named")
  }

  delete(p.clients, n.Name())
}

func (p *NullParent) List() []string {
  var result []string
  for k, _ := range p.clients {
    result = append(result, k)
  }
  return result
}

func TestClientName(t *testing.T) {
  client := NewClient(NewNullConn())

  name := "test"

  client.SetName(name)

  if client.Name() != name {
    t.Fatal("client returned wrong name", name)
  }
}

func TestClientParent(t *testing.T) {
  parent1 := NewNullParent()
  client := NewClient(NewNullConn())

  client.SetParent(parent1)

  parent2 := client.Parent()

  pValue1 := reflect.ValueOf(parent1)
  pValue2 := reflect.ValueOf(parent2)

  switch pValue2.Kind() {
  case reflect.Pointer:
    t.Log("parent2 is a pointer")
  default:
    t.Fatal("parent2 is NOT a pointer")
  }

  if !pValue1.Equal(pValue2) {
    t.Fatal("parent1 is NOT equal parent2")
  } else {
    t.Log("parent1 equal parent2")
  }
}

func TestClientSend(t *testing.T) {
  conn := NewNullConn()
  client := NewClient(conn)

  go client.Run(context.Background())
  time.Sleep(100*time.Millisecond)

  var d = []byte{1, 2, 3}

  client.Send() <- d
  result := <-conn.Messages()

  d = append(d, []byte{'\n'}...)

  if !reflect.DeepEqual(d, result) {
    t.Fatal("conn received different value than was sent")
  }
}

