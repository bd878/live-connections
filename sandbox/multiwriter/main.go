package main

import (
  "fmt"
  "io"
  "strings"
  "encoding/binary"
)

type A int

type B string

func (a *A) Write(p []byte) (int, error) {
  x := binary.LittleEndian.Uint16(p)
  *a = A(x)
  return len(p), nil
}

func (b *B) Write(p []byte) (int, error) {
  var sb strings.Builder
  n, err := sb.Write(p)
  *b = B(sb.String())
  return n, err
}

func main() {
  var a A
  var b B

  w := io.MultiWriter(&a, &b)
  w.Write([]byte{0x3e, 0x03})

  fmt.Println(a)
  fmt.Println(b)
}