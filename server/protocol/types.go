package protocol

type Encoder interface {
  Encode() []byte
}

type Decoder interface {
  Decode() (Encoder, error)
}

type Typed struct {
  messageType int8
}

type Coords struct {
  XPos float32
  YPos float32
}

type Text struct {
  text string
}

type Raw struct {
  size uint16
  data []byte
}

type List struct {
  items []string
}

type Identity struct {
  user string
}