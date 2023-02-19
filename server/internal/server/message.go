package server

import (
  "io"
  "encoding/binary"
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
  var err error
  if err = binary.Read(r, enc, &m.size); err != nil {
    meta.Log().Warn("binary.Read size err =", err)
    return 0, err
  }

  m.raw = make([]byte, m.size)
  if err = binary.Read(r, enc, m.raw); err != nil {
    meta.Log().Warn("binary.Read message err =", err)
    return 0, err
  }

  if m.size != uint16(len(m.raw)) {
    meta.Log().Warn("size != message size = ", m.size, len(m.raw))
    return 0, err
  }

  mr := bytes.NewReader(m.raw)
  if err = binary.Read(mr, enc, &m.messageType); err != nil {
    meta.Log().Warn("failed to read message type")
    return 0, err
  }

  return int64(m.size), nil
}

func (m *Message) Type() int8 {
  return m.messageType
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

func (m *Message) parseAuthMessage() error {
  mr := bytes.NewReader(m.raw)

  var err error

  var areaSize uint16
  if err := binary.Read(mr, enc, &areaSize); err != nil {
    return err
  }

  areaBytes := make([]byte, areaSize)
  if err = binary.Read(mr, enc, &areaBytes); err != nil {
    return err
  }

  var userSize uint16
  if err = binary.Read(mr, enc, &userSize); err != nil {
    return err
  }

  userBytes := make([]byte, userSize)
  if err = binary.Read(mr, enc, &userBytes); err != nil {
    return err
  }

  m.area = string(areaBytes)
  m.user = string(userBytes)

  return nil
}

func (m *Message) parseMouseMoveMessage() error {
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

func (m *Message) encodeTextMessage() []byte {
  message := []byte(m.text)

  result := new(bytes.Buffer)
  var size uint16 = uint16(len(message))
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

func (m *Message) encodeMouseMoveMessage() []byte {
  var err error

  typeBytes := new(bytes.Buffer)
  if err = binary.Write(typeBytes, enc, mouseMoveMessageType); err != nil {
    return nil
  }

  userBytes := new(bytes.Buffer)
  if err = binary.Write(userBytes, enc, []byte(m.user)); err != nil {
    return nil
  }

  userSizeBytes := new(bytes.Buffer)
  if err = binary.Write(userSizeBytes, enc, uint16(userBytes.Len())); err != nil {
    return nil
  }

  coordsBytes := new(bytes.Buffer)
  if err = binary.Write(coordsBytes, enc, m.XPos); err != nil {
    return nil
  }

  if err = binary.Write(coordsBytes, enc, m.YPos); err != nil {
    return nil
  }

  size := typeBytes.Len() + userSizeBytes.Len() + userBytes.Len() + coordsBytes.Len()
  sizeBytes := new(bytes.Buffer)
  if err = binary.Write(sizeBytes, enc, uint16(size)); err != nil {
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
  var err error

  typeBytes := new(bytes.Buffer)
  if err = binary.Write(typeBytes, enc, listClientsOnlineMessageType); err != nil {
    meta.Log().Warn("error writing message type =", err)
    return nil
  }

  usersBytes := new(bytes.Buffer)
  for _, user := range m.clients {
    var size uint16 = uint16(len(user))
    if err = binary.Write(usersBytes, enc, size); err != nil {
      meta.Log().Warn("error writing user name size =", err)
      return nil
    }

    if err = binary.Write(usersBytes, enc, []byte(user)); err != nil {
      meta.Log().Warn("error writing user name size =", err)
      return nil
    }
  }

  size := typeBytes.Len() + usersBytes.Len()
  sizeBytes := new(bytes.Buffer)
  if err = binary.Write(sizeBytes, enc, uint16(size)); err != nil {
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