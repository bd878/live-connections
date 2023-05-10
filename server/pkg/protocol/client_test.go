package protocol

import (
  "testing"
  "reflect"
  "io"
  "fmt"
)

type NullConn struct {}

func (c *NullConn) NextReader() (int, io.Reader, error) {
  return 0, nil, nil
}

func (c *NullConn) WriteMessage(mtype int, data []byte) error {
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

func TestClient(_ *testing.T) {
  parent1 := NewNullParent()
  client := NewClient(&NullConn{})

  client.SetName("test")
  client.SetParent(parent1)

  parent2 := client.Parent()

  pValue1 := reflect.ValueOf(parent1)
  pValue2 := reflect.ValueOf(parent2)
  switch pValue1.Kind() {
  case reflect.Pointer:
    fmt.Println("parent1 is pointer")
  default:
    fmt.Println("parent1 is NOT pointer")
  }

  switch pValue2.Kind() {
  case reflect.Pointer:
    fmt.Println("parent2 is pointer")
  default:
    fmt.Println("parent2 is NOT pointer")
  }

  if pValue1.Equal(pValue2) {
    fmt.Println("parent1 is equal parent2")
  }

  fmt.Println(client.Name())
}
