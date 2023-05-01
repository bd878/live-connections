package main

import "fmt"

type Text struct {
  Value string
}

type Printer interface {
  Print()
}

func (t *Text) Print() {
  fmt.Println(t.Value)
}

func MakePrinter() Printer {
  return &Text{}
}

func main() {
  p := MakePrinter()
  t, ok := p.(*Text)
  if !ok {
    fmt.Println("p is not Printer interface")
    return
  }

  t.Value = "aaa"
  t.Print()
}