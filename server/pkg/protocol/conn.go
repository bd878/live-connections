package protocol

import (
  "io"
  "time"

  ws "github.com/gorilla/websocket"
)

type WSConn struct {
  *ws.Conn
}

func NewConn(wsConn *ws.Conn) *WSConn {
  return &WSConn{wsConn}
}

func (p *WSConn) Init() {
  p.SetReadLimit(MaxPayloadSize)
  p.SetReadDeadline(time.Now().Add(PongWait))
  p.SetPongHandler(func(string) error { p.SetReadDeadline(time.Now().Add(PongWait)); return nil })
}

func (p *WSConn) NextReader() (int, io.Reader, error) {
  return p.NextReader()
}

func (p *WSConn) WriteMessage(messageType int, data []byte) error {
  p.SetWriteDeadline(time.Now().Add(WriteWait))

  return p.WriteMessage(messageType, data)
}

func (p *WSConn) Close() error {
  return p.Close()
}
