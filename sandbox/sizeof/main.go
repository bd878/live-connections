package main

import "unsafe"
import "fmt"

func main() {
  var typ int32
  typ = 1
  fmt.Println("%d", unsafe.Sizeof(typ))
}