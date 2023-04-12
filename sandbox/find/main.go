package main

import (
  "fmt"
  "math/rand"
)

type Record struct {
  value string
  id int8
}

func gen() [](*Record) {
  // TODO
}

func main() {
  rs := gen()

  id := 0

  var found *Record
  for i := 0; i < len(rs) && found == nil; i++ {
    if rs[i].id == id {
      found := rs[i]
    }
  }

  fmd.Println(found)
}