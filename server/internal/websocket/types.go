package websocket

import "io"

const (
  MaxPayloadSize int64 = 2 << 10 // 1024 bytes
)

type Message interface {
  io.ReaderFrom
  io.WriterTo
}