package messages

type Encoder interface {
  Encode() []byte
}

type Decoder interface {
  Decode() (Encoder, error)
}

type Typed struct {
  MessageType int8
}

type Coords struct {
  XPos float32
  YPos float32
}

type Text struct {
  Value string
}

type Raw struct {
  Size uint16
  Data []byte
}

type List struct {
  Items []string
}

type TextRecord struct {
  Text Text
  Title string
  ID int32
  UpdatedAt int32
  CreatedAt int32
}

type RecordsList struct {
  Items [](*TextRecord)
}

type Identity struct {
  User string
  Area string
}