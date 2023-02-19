package server

import (
  "io"
  "encoding/binary"
  "unsafe"
  "bytes"

  "github.com/teralion/live-connections/server/internal/meta"
)

var enc = binary.LittleEndian

const (
  authMessageType int8 = 1
  mouseMoveMessageType int8 = 2
  listClientsOnlineMessageType int8 = 3
  authOkMessageType int8 = 4
  initMouseCoordsMessageType int8 = 5
  textMessageType int8 = 6
)

type Message struct {
  size uint16
  messageType int8
  raw []byte

  area string
  user string

  XPos float32
  YPos float32

  text string

  clients []string
}

func NewMessage() *Message {
  return &Message{}
}

// size + type + raw
func (m *Message) ReadFrom(r io.Reader) (int64, error) {
  if err := binary.Read(r, enc, &m.size); err != nil {
    meta.Log().Warn("binary.Read size err =", err)
    return 0, err
  }

  if err := binary.Read(r, enc, &m.messageType); err != nil {
    meta.Log().Warn("failed to read message type")
    return 0, err
  }

  m.size = m.size - uint16(unsafe.Sizeof(m.messageType)) // size = type + raw

  m.raw = make([]byte, m.size)
  if err := binary.Read(r, enc, &m.raw); err != nil {
    meta.Log().Warn("binary.Read message err =", err)
    return 0, err
  }

  return int64(m.size), nil
}

func (m *Message) Type() int8 {
  return m.messageType
}

func (m *Message) SetArea(area string) {
  m.area = area
}

func (m *Message) SetUser(user string) {
  m.user = user
}

func (m *Message) Decode() error {
  switch (m.Type()) {
  case authMessageType:
    return m.parseAuthMessage()
  case mouseMoveMessageType:
    return m.parseMouseMoveMessage()
  }
  return nil
}

func (m *Message) Encode() []byte {
  switch (m.Type()) {
  case authMessageType:
  case textMessageType:
    return m.encodeTextMessage()
  case mouseMoveMessageType:
    return m.encodeMouseMoveMessage()
  case listClientsOnlineMessageType:
    return m.encodeClientsOnlineMessage()
  }
  return nil
}

// areaSize + areaBytes + userSize + userBytes
func (m *Message) parseAuthMessage() error {
  meta.Log().Debug("parse auth message")

  mr := bytes.NewReader(m.raw)

  var areaSize uint16
  if err := binary.Read(mr, enc, &areaSize); err != nil {
    meta.Log().Warn("failed to read areaSize")
    return err
  }

  areaBytes := make([]byte, areaSize)
  if err := binary.Read(mr, enc, &areaBytes); err != nil {
    meta.Log().Warn("failed to read area")
    return err
  }

  var userSize uint16
  if err := binary.Read(mr, enc, &userSize); err != nil {
    meta.Log().Warn("failed to read userSize")
    return err
  }

  userBytes := make([]byte, userSize)
  if err := binary.Read(mr, enc, &userBytes); err != nil {
    meta.Log().Warn("failed to read user")
    return err
  }

  m.area = string(areaBytes)
  m.user = string(userBytes)

  return nil
}

// XPos + YPos
func (m *Message) parseMouseMoveMessage() error {
  meta.Log().Debug("parse mouse move message")

  mr := bytes.NewReader(m.raw)

  if err := binary.Read(mr, enc, &m.XPos); err != nil {
    meta.Log().Warn("failed to read xPos =",  err)
    return err
  }

  if err := binary.Read(mr, enc, &m.YPos); err != nil {
    meta.Log().Warn("failed to read YPos =",  err)
    return err
  }

  return nil
}

// totalSize + type + text message
func (m *Message) encodeTextMessage() []byte {
  meta.Log().Debug("encode text message")

  message := []byte(m.text)

  result := new(bytes.Buffer)
  size := uint16(len(message))
  if err := binary.Write(result, enc, size); err != nil {
    meta.Log().Warn("error writing text message size =", err)
    return nil
  }

  if err := binary.Write(result, enc, m.messageType); err != nil {
    meta.Log().Warn("error writing text message type =", err)
    return nil
  }

  if err := binary.Write(result, enc, message); err != nil {
    meta.Log().Warn("error writing text message =", err)
    return nil
  }

  return result.Bytes()
}

// totalSize + type + userSize + userBytes + XPos + YPos
func (m *Message) encodeMouseMoveMessage() []byte {
  meta.Log().Debug("encode mouse move message")

  typeBytes := new(bytes.Buffer)
  if err := binary.Write(typeBytes, enc, mouseMoveMessageType); err != nil {
    meta.Log().Warn("error writing mouse move type")
    return nil
  }

  userBytes := new(bytes.Buffer)
  if err := binary.Write(userBytes, enc, []byte(m.user)); err != nil {
    meta.Log().Warn("error writing user size")
    return nil
  }

  userSizeBytes := new(bytes.Buffer)
  if err := binary.Write(userSizeBytes, enc, uint16(userBytes.Len())); err != nil {
    meta.Log().Warn("error writing user")
    return nil
  }

  coordsBytes := new(bytes.Buffer)
  if err := binary.Write(coordsBytes, enc, m.XPos); err != nil {
    meta.Log().Warn("error writing y coord")
    return nil
  }

  if err := binary.Write(coordsBytes, enc, m.YPos); err != nil {
    meta.Log().Warn("error writing y coord")
    return nil
  }

  size := typeBytes.Len() + userSizeBytes.Len() + userBytes.Len() + coordsBytes.Len()
  sizeBytes := new(bytes.Buffer)
  if err := binary.Write(sizeBytes, enc, uint16(size)); err != nil {
    meta.Log().Warn("error writing size")
    return nil
  }

  return bytes.Join(
    [][]byte{
      sizeBytes.Bytes(),
      typeBytes.Bytes(),
      userSizeBytes.Bytes(),
      userBytes.Bytes(),
      coordsBytes.Bytes(),
    },
    []byte{},
  )
}

func EncodeClientsOnline(clients []string) []byte {
  m := &Message{clients: clients}
  return m.encodeClientsOnlineMessage()
}

func (m *Message) encodeClientsOnlineMessage() []byte {
  meta.Log().Debug("encode clients online message")

  typeBytes := new(bytes.Buffer)
  if err := binary.Write(typeBytes, enc, listClientsOnlineMessageType); err != nil {
    meta.Log().Warn("error writing message type =", err)
    return nil
  }

  usersBytes := new(bytes.Buffer)
  for _, user := range m.clients {
    var size uint16 = uint16(len(user))
    if err := binary.Write(usersBytes, enc, size); err != nil {
      meta.Log().Warn("error writing user name size =", err)
      return nil
    }

    if err := binary.Write(usersBytes, enc, []byte(user)); err != nil {
      meta.Log().Warn("error writing user name size =", err)
      return nil
    }
  }

  size := typeBytes.Len() + usersBytes.Len()
  sizeBytes := new(bytes.Buffer)
  if err := binary.Write(sizeBytes, enc, uint16(size)); err != nil {
    meta.Log().Warn("error writing total size =", err)
    return nil
  }

  return bytes.Join(
    [][]byte{
      sizeBytes.Bytes(),
      typeBytes.Bytes(),
      usersBytes.Bytes(),
    },
    []byte{},
  )
}