package protocol

import (
  "io"
  "time"

  ws "github.com/gorilla/websocket"
)

type ProtocolConn struct {
  *ws.Conn
}

func NewConn(wsConn *ws.Conn) *ProtocolConn {
  return &ProtocolConn{wsConn}
}

func (p *ProtocolConn) Init() {
  p.SetReadLimit(MaxPayloadSize)
  p.SetReadDeadline(time.Now().Add(PongWait))
  p.SetPongHandler(func(string) error { p.SetReadDeadline(time.Now().Add(PongWait)); return nil })
}

func (p *ProtocolConn) NextReader() (int, io.Reader, error) {
  return p.NextReader()
}

func (p *ProtocolConn) WriteMessage(messageType int, data []byte) error {
  p.SetWriteDeadline(time.Now().Add(WriteWait))

  return p.WriteMessage(messageType, data)
}
