package main

import (
  "log"
  "unicode/utf8"
)

func main() {
  barr := []byte{0xe7, 0xa5, 0xa8}
  log.Println("valid =", utf8.Valid(barr))
  log.Println("str =", string(barr))
}