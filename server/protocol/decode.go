package protocol

import (
  "io"
  "encoding/binary"
  "unsafe"
  "bytes"
  "errors"

  "github.com/teralion/live-connections/meta"
)

type RawMessage struct {
  Typed
  Raw
}

// size + type + raw
func ReadFrom(r io.Reader) (Decoder, error) {
  raw := Raw{}

  if err := binary.Read(r, enc, &raw.size); err != nil {
    meta.Log().Warn("binary.Read size err =", err)
    return nil, err
  }

  typed := Typed{}
  if err := binary.Read(r, enc, &typed.messageType); err != nil {
    meta.Log().Warn("failed to read message type")
    return nil, err
  }

  raw.size = raw.size - uint16(unsafe.Sizeof(typed.messageType)) // size = type + raw

  if err := binary.Read(r, enc, &raw.data); err != nil {
    meta.Log().Warn("binary.Read message err =", err)
    return nil, err
  }

  return &RawMessage{typed, raw}, nil
}

func (m *RawMessage) Decode() (Encoder, error) {
  switch (m.messageType) {
  case auth:
    return m.decodeAuth()
  case mouseMove:
    return m.decodeMouseMove()
  case squareMove:
    return m.decodeSquareMove()
  case text:
    return m.decodeText()
  default:
    return nil, errors.New("undefined message type")
  }
}

// areaSize + areaBytes + userSize + userBytes
func (m *RawMessage) decodeAuth() (Encoder, error) {
  meta.Log().Debug("decode auth")

  mr := bytes.NewReader(m.data)

  var areaSize uint16
  if err := binary.Read(mr, enc, &areaSize); err != nil {
    meta.Log().Warn("failed to read areaSize")
    return nil, err
  }

  areaBytes := make([]byte, areaSize)
  if err := binary.Read(mr, enc, &areaBytes); err != nil {
    meta.Log().Warn("failed to read area")
    return nil, err
  }

  var userSize uint16
  if err := binary.Read(mr, enc, &userSize); err != nil {
    meta.Log().Warn("failed to read userSize")
    return nil, err
  }

  userBytes := make([]byte, userSize)
  if err := binary.Read(mr, enc, &userBytes); err != nil {
    meta.Log().Warn("failed to read user")
    return nil, err
  }

  result := AuthMessage{}
  meta.Log().Debug("area =", string(areaBytes))
  result.SetUser(string(userBytes))

  return &result, nil
}

func (m *RawMessage) decodeMouseMove() (Encoder, error) {
  meta.Log().Debug("decode mouse move")

  msg := MouseMoveMessage{}
  var err error
  msg.Coords, err = m.decodeCoords()
  if err != nil {
    return nil, err
  }
  return &msg, nil
}

func (m *RawMessage) decodeSquareMove() (Encoder, error) {
  meta.Log().Debug("square move")

  msg := SquareMoveMessage{}
  var err error
  msg.Coords, err = m.decodeCoords()
  if err != nil {
    return nil, err
  }
  return &msg, nil
}

// XPos + YPos
func (m *RawMessage) decodeCoords() (Coords, error) {
  meta.Log().Debug("decode coords")

  mr := bytes.NewReader(m.data)

  result := Coords{}

  if err := binary.Read(mr, enc, &result.XPos); err != nil {
    meta.Log().Warn("failed to read XPos =",  err)
    return result, err
  }

  if err := binary.Read(mr, enc, &result.YPos); err != nil {
    meta.Log().Warn("failed to read YPos =",  err)
    return result, err
  }

  return result, nil
}

// textSize + text message
func (m *RawMessage) decodeText() (Encoder, error) {
  meta.Log().Debug("parse text message")

  mr := bytes.NewReader(m.data)

  msg := TextMessage{}

  var textSize uint16
  if err := binary.Read(mr, enc, &textSize); err != nil {
    meta.Log().Warn("failed to read textSize")
    return nil, err
  }

  textBytes := make([]byte, textSize)
  if err := binary.Read(mr, enc, &textBytes); err != nil {
    meta.Log().Warn("failed to read textBytes =",  err)
    return nil, err
  }

  msg.SetText(string(textBytes))

  return &msg, nil
}