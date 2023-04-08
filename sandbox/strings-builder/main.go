package main

import (
  "log"
  "fmt"
  "unicode/utf8"
  "strings"
  "time"
)

func getFileName() {
  var b strings.Builder
  fmt.Fprintf(&b, "%d.state", int32(time.Now().Unix()))
  fmt.Println(b.String())
}

func joinBytes() {
  barr := []byte{0xe7, 0xa5, 0xa8}
  log.Println("valid =", utf8.Valid(barr))
  log.Println("str =", string(barr))
}

func main() {
  getFileName()
}