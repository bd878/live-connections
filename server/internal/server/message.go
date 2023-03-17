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
  squareInitMessageType int8 = 7
  squareMoveMessageType int8 = 8
  textInputMessageType int8 = 9
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

type ClientCoords struct {
  name string
  XPos float32
  YPos float32
}

type ClientText struct {
  name string
  text string
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
    panic(err)
    return 0, err
  }

  return int64(m.size), nil
}

func (m *Message) Type() int8 {
  return m.messageType
}

func (m *Message) SetType(typ int8) {
  m.messageType = typ
}

func (m *Message) SetArea(area string) {
  m.area = area
}

func (m *Message) SetUser(user string) {
  m.user = user
}

func (m *Message) SetClients(clients []string) {
  m.clients = clients
}

func (m *Message) SetX(XPos float32) {
  m.XPos = XPos
}

func (m *Message) SetY(YPos float32) {
  m.YPos = YPos
}

func (m *Message) Decode() error {
  switch (m.Type()) {
  case authMessageType:
    return m.parseAuthMessage()
  case mouseMoveMessageType:
    return m.parseCoordsMessage()
  case squareMoveMessageType:
    return m.parseCoordsMessage()
  case textInputMessageType:
    return m.parseTextMessage()
  }
  return nil
}

func (m *Message) Encode() []byte {
  switch (m.Type()) {
  case authMessageType:
    return m.encodeAuthMessage()
  case textInputMessageType:
    return m.encodeTextInputMessage()
  case mouseMoveMessageType:
    return m.encodeCoordsMessage()
  case squareInitMessageType:
    return m.encodeCoordsMessage()
  case squareMoveMessageType:
    return m.encodeCoordsMessage()
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
func (m *Message) parseCoordsMessage() error {
  meta.Log().Debug("parse coords message")

  mr := bytes.NewReader(m.raw)

  if err := binary.Read(mr, enc, &m.XPos); err != nil {
    meta.Log().Warn("failed to read XPos =",  err)
    return err
  }

  if err := binary.Read(mr, enc, &m.YPos); err != nil {
    meta.Log().Warn("failed to read YPos =",  err)
    return err
  }

  return nil
}

// textSize + text message
func (m *Message) parseTextMessage() error {
  meta.Log().Debug("parse text message")

  mr := bytes.NewReader(m.raw)

  var textSize uint16
  if err := binary.Read(mr, enc, &textSize); err != nil {
    meta.Log().Warn("failed to read textSize")
    return err
  }

  textBytes := make([]byte, textSize)
  if err := binary.Read(mr, enc, &textBytes); err != nil {
    meta.Log().Warn("failed to read textBytes =",  err)
    return err
  }

  m.text = string(textBytes)

  return nil
}

// totalSize + type + text message
func (m *Message) encodeAuthMessage() []byte {
  meta.Log().Debug("encode auth message")

  message := []byte(m.text)

  result := new(bytes.Buffer)
  size := uint16(len(message))
  if err := binary.Write(result, enc, size); err != nil {
    meta.Log().Warn("error writing auth message size =", err)
    return nil
  }

  if err := binary.Write(result, enc, m.messageType); err != nil {
    meta.Log().Warn("error writing auth message type =", err)
    return nil
  }

  if err := binary.Write(result, enc, message); err != nil {
    meta.Log().Warn("error writing auth message =", err)
    return nil
  }

  return result.Bytes()
}

// TODO: retrieve common encodings (type, user ...etc) in interface
// totalSize + type + userSize + userBytes + textSize + textBytes
func (m *Message) encodeTextInputMessage() []byte {
  meta.Log().Debug("encode text input message")

  typeBytes := new(bytes.Buffer)
  if err := binary.Write(typeBytes, enc, m.messageType); err != nil {
    meta.Log().Warn("error writing message type type")
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

  textBytes := new(bytes.Buffer)
  if err := binary.Write(textBytes, enc, []byte(m.text)); err != nil {
    meta.Log().Warn("error writing text size")
    return nil
  }

  textSizeBytes := new(bytes.Buffer)
  if err := binary.Write(textSizeBytes, enc, uint16(textBytes.Len())); err != nil {
    meta.Log().Warn("error writing user")
    return nil
  }

  size := typeBytes.Len() + userSizeBytes.Len() + userBytes.Len() + textSizeBytes.Len() + textBytes.Len()
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
      textSizeBytes.Bytes(),
      textBytes.Bytes(),
    },
    []byte{},
  )
}

// totalSize + type + userSize + userBytes + XPos + YPos
func (m *Message) encodeCoordsMessage() []byte {
  meta.Log().Debug("encode coords message")

  typeBytes := new(bytes.Buffer)
  if err := binary.Write(typeBytes, enc, m.messageType); err != nil {
    meta.Log().Warn("error writing coords type")
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
  m := &Message{messageType: listClientsOnlineMessageType, clients: clients}
  return m.Encode()
}

func EncodeSquareInit(c *ClientCoords) []byte {
  m := &Message{
    messageType: squareInitMessageType,
    user: c.name,
    XPos: c.XPos,
    YPos: c.YPos,
  }
  return m.Encode()
}

func EncodeTextInputInit(c *ClientText) []byte {
  m := &Message{
    messageType: textInputMessageType,
    user: c.name,
    text: c.text,
  }
  return m.Encode()
}

// totalSize + type + usersCount + []{userSize + userBytes}
func (m *Message) encodeClientsOnlineMessage() []byte {
  meta.Log().Debug("encode clients online message")

  typeBytes := new(bytes.Buffer)
  if err := binary.Write(typeBytes, enc, m.messageType); err != nil {
    meta.Log().Warn("error writing message type =", err)
    return nil
  }

  countBytes := new(bytes.Buffer)
  if err := binary.Write(countBytes, enc, uint16(len(m.clients))); err != nil {
    meta.Log().Warn("error writing users count =", err)
    return nil
  }

  usersBytes := new(bytes.Buffer)
  for _, user := range m.clients {
    size := uint16(len(user))
    if err := binary.Write(usersBytes, enc, size); err != nil {
      meta.Log().Warn("error writing user name size =", err)
      return nil
    }

    if err := binary.Write(usersBytes, enc, []byte(user)); err != nil {
      meta.Log().Warn("error writing user name size =", err)
      return nil
    }
  }

  size := typeBytes.Len() + countBytes.Len() + usersBytes.Len()
  sizeBytes := new(bytes.Buffer)
  if err := binary.Write(sizeBytes, enc, uint16(size)); err != nil {
    meta.Log().Warn("error writing total size =", err)
    return nil
  }

  return bytes.Join(
    [][]byte{
      sizeBytes.Bytes(),
      typeBytes.Bytes(),
      countBytes.Bytes(),
      usersBytes.Bytes(),
    },
    []byte{},
  )
}