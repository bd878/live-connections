package messages

import (
  "io"
  "encoding/binary"
  "unsafe"
  "bytes"
  "errors"

  "github.com/bd878/live-connections/meta"
)

type RawMessage struct {
  Typed
  Raw
}

// size + type + raw
func ReadFrom(r io.Reader) (Decoder, error) {
  raw := Raw{}

  if err := binary.Read(r, enc, &raw.Size); err != nil {
    meta.Log().Warn("binary.Read size err =", err)
    return nil, err
  }

  typed := Typed{}
  if err := binary.Read(r, enc, &typed.MessageType); err != nil {
    meta.Log().Warn("failed to read message type")
    return nil, err
  }

  raw.Size = raw.Size - uint16(unsafe.Sizeof(typed.MessageType)) // size = type + raw

  raw.Data = make([]byte, raw.Size)
  if err := binary.Read(r, enc, &raw.Data); err != nil {
    meta.Log().Warn("binary.Read message err =", err)
    return nil, err
  }

  return &RawMessage{Typed: typed, Raw: raw}, nil
}

func (m *RawMessage) Decode() (Encoder, error) {
  switch (m.MessageType) {
  case auth:
    return m.decodeAuth()
  case mouseMove:
    return m.decodeMouseMove()
  case squareMove:
    return m.decodeSquareMove()
  case text:
    return m.decodeText()
  case addRecord:
    return m.decodeAddRecord()
  case selectRecord:
    return m.decodeSelectRecord()
  default:
    return nil, errors.New("undefined message type")
  }
}

// areaSize + areaBytes + userSize + userBytes
func (m *RawMessage) decodeAuth() (Encoder, error) {
  meta.Log().Debug("decode auth")

  mr := bytes.NewReader(m.Raw.Data)

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

  typed := Typed{MessageType: auth}
  result := AuthMessage{Typed: typed}
  result.SetArea(string(areaBytes))
  result.SetUser(string(userBytes))

  return &result, nil
}

func (m *RawMessage) decodeMouseMove() (Encoder, error) {
  meta.Log().Debug("decode mouse move")

  typed := Typed{MessageType: mouseMove}
  msg := MouseMoveMessage{Typed: typed}
  var err error
  msg.Coords, err = m.decodeCoords()
  if err != nil {
    return nil, err
  }
  return &msg, nil
}

func (m *RawMessage) decodeSquareMove() (Encoder, error) {
  meta.Log().Debug("decode square move")

  typed := Typed{MessageType: squareMove}
  msg := SquareMoveMessage{Typed: typed}
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

  mr := bytes.NewReader(m.Data)

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
  meta.Log().Debug("decode text message")

  mr := bytes.NewReader(m.Data)

  typed := Typed{MessageType: text}
  result := TextMessage{Typed: typed}

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

  result.SetText(string(textBytes))

  return &result, nil
}

func (m *RawMessage) decodeAddRecord() (Encoder, error) {
  meta.Log().Debug("decode add record message")

  typed := Typed{MessageType: addRecord}
  result := AddRecordMessage{Typed: typed}

  return &result, nil
}

// createdAtBytes
func (m *RawMessage) decodeSelectRecord() (Encoder, error) {
  meta.Log().Debug("decode select record message")

  return nil, nil
}