package messages

import (
  "encoding/binary"
  "bytes"

  "github.com/bd878/live-connections/meta"
)

func (i *Identity) SetUser(name string) {
  i.User = name
}

func (i *Identity) SetArea(area string) {
  i.Area = area
}

func (t *Typed) SetType(typ int8) {
  t.MessageType = typ
}

func (c *Coords) SetX(XPos float32) {
  c.XPos = XPos
}

func (c *Coords) SetY(YPos float32) {
  c.YPos = YPos
}

func (l *List) SetItems(items []string) {
  l.Items = items
}

func (t *Text) SetText(str string) {
  t.Str = str
}

type AuthMessage struct {
  Typed
  Identity
  Text
}

// totalSize + type + text message
func (m *AuthMessage) Encode() []byte {
  typeBytes := m.Typed.Encode()
  textBytes := m.Text.Encode()

  sizeBytes := encodeSize(len(typeBytes) + len(textBytes))

  return bytes.Join(
    [][]byte{
      sizeBytes,
      typeBytes,
      textBytes,
    },
    []byte{},
  )
}

type SquareMoveMessage struct {
  Typed
  Coords
  Identity
}

func NewSquareMoveMessage(user string, XPos, YPos float32) *SquareMoveMessage {
  message := &SquareMoveMessage{
    Typed{MessageType: squareMove},
    Coords{XPos: XPos, YPos: YPos},
    Identity{User: user},
  }
  return message
}

// totalSize + type + userSize + userBytes + XPos + YPos
func (m *SquareMoveMessage) Encode() []byte {
  typeBytes := m.Typed.Encode()
  userBytes := m.Identity.Encode()
  coordsBytes := m.Coords.Encode()

  sizeBytes := encodeSize(len(typeBytes) + len(userBytes) + len(coordsBytes))

  return bytes.Join(
    [][]byte{
      sizeBytes,
      typeBytes,
      userBytes,
      coordsBytes,
    },
    []byte{},
  )
}

type SquareInitMessage struct {
  Typed
  Coords
  Identity
}

func NewSquareInitMessage(user string, XPos, YPos float32) *SquareInitMessage {
  message := &SquareInitMessage{
    Typed{MessageType: squareInit},
    Coords{XPos: XPos, YPos: YPos},
    Identity{User: user},
  }
  return message
}

// totalSize + type + userSize + userBytes + XPos + YPos
func (m *SquareInitMessage) Encode() []byte {
  typeBytes := m.Typed.Encode()
  userBytes := m.Identity.Encode()
  coordsBytes := m.Coords.Encode()

  sizeBytes := encodeSize(len(typeBytes) + len(userBytes) + len(coordsBytes))

  return bytes.Join(
    [][]byte{
      sizeBytes,
      typeBytes,
      userBytes,
      coordsBytes,
    },
    []byte{},
  )
}

type MouseMoveMessage struct {
  Typed
  Coords
  Identity
}

func NewMouseMoveMessage(user string, XPos, YPos float32) *MouseMoveMessage {
  message := &MouseMoveMessage{
    Typed{MessageType: mouseMove},
    Coords{XPos: XPos, YPos: YPos},
    Identity{User: user},
  }
  return message
}

// totalSize + type + userSize + userBytes + XPos + YPos
func (m *MouseMoveMessage) Encode() []byte {
  typeBytes := m.Typed.Encode()
  userBytes := m.Identity.Encode()
  coordsBytes := m.Coords.Encode()

  sizeBytes := encodeSize(len(typeBytes) + len(userBytes) + len(coordsBytes))

  return bytes.Join(
    [][]byte{
      sizeBytes,
      typeBytes,
      userBytes,
      coordsBytes,
    },
    []byte{},
  )
}

type TitlesListMessage struct {
  Typed
  RecordsList
  Identity
}

func NewTitlesListMessage(user string, items [](*Record)) *TitlesListMessage {
  message := &TitlesListMessage{
    Typed{MessageType: titlesList},
    RecordsList{Items: items},
    Identity{User: user},
  }
  return message
}

func (m *TitlesListMessage) Encode() []byte {
  typeBytes := m.Typed.Encode()
  userBytes := m.Identity.Encode()
  recordsBytes := m.RecordsList.Encode()

  sizeBytes := encodeSize(len(typeBytes) + len(userBytes) + len(recordsBytes))

  return bytes.Join(
    [][]byte{
      sizeBytes,
      typeBytes,
      userBytes,
      recordsBytes,
    },
    []byte{},
  )
}

type ClientsOnlineMessage struct {
  Typed
  List
}

func NewClientsOnlineMessage(items []string) *ClientsOnlineMessage {
  message := &ClientsOnlineMessage{
    Typed{MessageType: clientsOnline},
    List{Items: items},
  }
  return message
}

// totalSize + type + usersCount + []{userSize + userBytes}
func (m *ClientsOnlineMessage) Encode() []byte {
  typeBytes := m.Typed.Encode()
  itemsBytes := m.List.Encode()

  sizeBytes := encodeSize(len(typeBytes) + len(itemsBytes))

  return bytes.Join(
    [][]byte{
      sizeBytes,
      typeBytes,
      itemsBytes,
    },
    []byte{},
  )
}

type AddRecordMessage struct {
  Typed
}

// totalSize + type
func (m *AddRecordMessage) Encode() []byte {
  typeBytes := m.Typed.Encode()

  sizeBytes := encodeSize(len(typeBytes))

  return bytes.Join(
    [][]byte{
      sizeBytes,
      typeBytes,
    },
    []byte{},
  )
}

type SelectRecordMessage struct {
  Typed
  Identity
  CreatedAt int32 // TODO: change on id
}

// totalSize + type + userSize + userBytes + createdAtBytes
func (m *SelectRecordMessage) Encode() []byte {
  typeBytes := m.Typed.Encode()
  userBytes := m.Identity.Encode()
  timeBytes := encodeUnixTime(m.CreatedAt)

  sizeBytes := encodeSize(len(typeBytes) + len(userBytes) + len(timeBytes))

  return bytes.Join(
    [][]byte{
      sizeBytes,
      typeBytes,
      userBytes,
      timeBytes,
    },
    []byte{},
  )
}

type TextMessage struct {
  Typed
  Identity
  Text
}

func NewTextMessage(user, str string) *TextMessage {
  message := &TextMessage{
    Typed{MessageType: text},
    Identity{User: user},
    Text{Str: str},
  }
  return message
}

// totalSize + type + userSize + userBytes + textSize + textBytes
func (m *TextMessage) Encode() []byte {
  typeBytes := m.Typed.Encode()
  userBytes := m.Identity.Encode()
  textBytes := m.Text.Encode()

  sizeBytes := encodeSize(len(typeBytes) + len(userBytes) + len(textBytes))

  return bytes.Join(
    [][]byte{
      sizeBytes,
      typeBytes,
      userBytes,
      textBytes,
    },
    []byte{},
  )
}

func NewCoords(XPos, YPos float32) *Coords {
  return &Coords{
    XPos: XPos,
    YPos: YPos,
  }
}

// XPos + YPos
func (m *Coords) Encode() []byte {
  coordsBytes := new(bytes.Buffer)
  if err := binary.Write(coordsBytes, enc, m.XPos); err != nil {
    meta.Log().Warn("error writing y coord")
    return nil
  }

  if err := binary.Write(coordsBytes, enc, m.YPos); err != nil {
    meta.Log().Warn("error writing y coord")
    return nil
  }

  return coordsBytes.Bytes()
}

// user size + user text
func (m *Identity) Encode() []byte {
  userBytes := new(bytes.Buffer)
  if err := binary.Write(userBytes, enc, []byte(m.User)); err != nil {
    meta.Log().Warn("error writing user size")
    return nil
  }

  userSizeBytes := new(bytes.Buffer)
  if err := binary.Write(userSizeBytes, enc, uint16(userBytes.Len())); err != nil {
    meta.Log().Warn("error writing user")
    return nil
  }

  return bytes.Join(
    [][]byte{
      userSizeBytes.Bytes(),
      userBytes.Bytes(),
    },
    []byte{},
  )
}

func NewText(str string) *Text {
  return &Text{
    Str: str,
  }
}

// text size + text
func (m *Text) Encode() []byte {
  textBytes := new(bytes.Buffer)
  if err := binary.Write(textBytes, enc, []byte(m.Str)); err != nil {
    meta.Log().Warn("error writing text size")
    return nil
  }

  textSizeBytes := new(bytes.Buffer)
  if err := binary.Write(textSizeBytes, enc, uint16(textBytes.Len())); err != nil {
    meta.Log().Warn("error writing user")
    return nil
  }

  return bytes.Join(
    [][]byte{
      textSizeBytes.Bytes(),
      textBytes.Bytes(),
    },
    []byte{},
  )
}

func (m *Typed) Encode() []byte {
  typeBytes := new(bytes.Buffer)
  if err := binary.Write(typeBytes, enc, m.MessageType); err != nil {
    meta.Log().Warn("error writing coords type")
    return nil
  }

  return typeBytes.Bytes()
}

// itemsCount + []{itemSize + itemBytes}
func (m *List) Encode() []byte {
  countBytes := new(bytes.Buffer)
  if err := binary.Write(countBytes, enc, uint16(len(m.Items))); err != nil {
    meta.Log().Warn("error writing items count =", err)
    return nil
  }

  itemsBytes := new(bytes.Buffer)
  for _, item := range m.Items {
    size := uint16(len(item))
    if err := binary.Write(itemsBytes, enc, size); err != nil {
      meta.Log().Warn("error writing item name size =", err)
      return nil
    }

    if err := binary.Write(itemsBytes, enc, []byte(item)); err != nil {
      meta.Log().Warn("error writing item name size =", err)
      return nil
    }
  }

  return bytes.Join(
    [][]byte{
      countBytes.Bytes(),
      itemsBytes.Bytes(),
    },
    []byte{},
  )
}

func NewRecordsList(items [](*Record)) *RecordsList {
  message := &RecordsList{
    Items: items,
  }
  return message
}

// recordsCount + []{itemSize + recordBytes}
func (m *RecordsList) Encode() []byte {
  countBytes := new(bytes.Buffer)
  if err := binary.Write(countBytes, enc, uint16(len(m.Items))); err != nil {
    meta.Log().Warn("error writing records count =", err)
    return nil
  }

  itemsBytes := new(bytes.Buffer)
  for _, item := range m.Items {
    rb := item.Encode()

    size := uint16(len(rb))
    if err := binary.Write(itemsBytes, enc, size); err != nil {
      meta.Log().Warn("error writing record name size =", err)
      return nil
    }

    if err := binary.Write(itemsBytes, enc, rb); err != nil {
      meta.Log().Warn("error writing record name size =", err)
      return nil
    }
  }

  return bytes.Join(
    [][]byte{
      countBytes.Bytes(),
      itemsBytes.Bytes(),
    },
    []byte{},
  )
}

// valueSize + valueBytes + idBytes + updatedAtBytes + createdAtBytes
func (m *Record) Encode() []byte {
  valueBytes := new(bytes.Buffer)
  if err := binary.Write(valueBytes, enc, []byte(m.Value)); err != nil {
    meta.Log().Warn("error writing record value")
    return nil
  }

  valueSize := new(bytes.Buffer)
  if err := binary.Write(valueSize, enc, uint16(valueBytes.Len())); err != nil {
    meta.Log().Warn("error writing record value size =", err)
    return nil
  }

  idBytes := new(bytes.Buffer)
  if err := binary.Write(idBytes, enc, m.ID); err != nil {
    meta.Log().Warn("error writing id value =", err)
    return nil
  }

  updatedAtBytes := encodeUnixTime(m.UpdatedAt)
  createdAtBytes := encodeUnixTime(m.CreatedAt)

  return bytes.Join(
    [][]byte{
      valueSize.Bytes(),
      valueBytes.Bytes(),
      idBytes.Bytes(),
      updatedAtBytes,
      createdAtBytes,
    },
    []byte{},
  )
}

func encodeUnixTime(t int32) []byte {
  timeBytes := new(bytes.Buffer)
  if err := binary.Write(timeBytes, enc, t); err != nil {
    meta.Log().Warn("error writing unix time value")
    return nil
  }
  return timeBytes.Bytes()
}

func encodeSize(size int) []byte {
  sizeBytes := new(bytes.Buffer)
  if err := binary.Write(sizeBytes, enc, uint16(size)); err != nil {
    meta.Log().Warn("error writing size")
    return nil
  }
  return sizeBytes.Bytes()
}