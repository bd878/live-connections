package protocol

import (
  "io"
  "time"

  ws "github.com/gorilla/websocket"
)

type WSConn struct {
  conn *ws.Conn
}

func NewConn(wsConn *ws.Conn) *WSConn {
  return &WSConn{conn: wsConn}
}

func (p *WSConn) Init() {
  p.conn.SetReadLimit(MaxPayloadSize)
  p.conn.SetReadDeadline(time.Now().Add(PongWait))
  p.conn.SetPongHandler(func(string) error { p.conn.SetReadDeadline(time.Now().Add(PongWait)); return nil })
}

func (p *WSConn) NextReader() (int, io.Reader, error) {
  return p.conn.NextReader()
}

func (p *WSConn) WriteMessage(messageType int, data []byte) error {
  p.conn.SetWriteDeadline(time.Now().Add(WriteWait))

  return p.conn.WriteMessage(messageType, data)
}

func (p *WSConn) Close() error {
  return p.conn.Close()
}
