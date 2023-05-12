package mock

import "io"

type Conn struct {
  messages chan []byte
}

func NewConn() *Conn {
  return &Conn{
    messages: make(chan []byte),
  }
}

func (c *Conn) Messages() chan []byte {
  return c.messages
}

func (c *Conn) NextReader() (int, io.Reader, error) {
  ch := make(chan struct{})
  <-ch

  return 0, nil, nil
}

func (c *Conn) WriteMessage(_ int, data []byte) error {
  c.Messages() <- data
  return nil
}

func (c *Conn) Close() error {
  return nil
}