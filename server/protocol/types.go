package protocol

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
  Str string
}

type Raw struct {
  Size uint16
  Data []byte
}

type List struct {
  Items []string
}

type Identity struct {
  User string
}