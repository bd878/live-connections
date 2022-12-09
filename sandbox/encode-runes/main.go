package main

import (
  "math/rand"
  "time"
  "unicode/utf8"
  "fmt"
)

func main() {
  rand.Seed(time.Now().UnixNano())

  b := make([]byte, 1<<4)
  n, err := rand.Read(b)
  if err != nil {
    fmt.Println("error: %v\n", err)
  }
  fmt.Println("read size: %d\n", n)
  fmt.Println("bytes =", len(b))
  fmt.Println("runes =", utf8.RuneCount(b))

  fmt.Println("bytes:")
  for i := 0; i < len(b); i++ {
    fmt.Printf("%U ", b[i])
  }
  fmt.Println()

  for len(b) > 0 {
    r, size := utf8.DecodeRune(b)
    fmt.Printf("%c %v\n", r, size)

    b = b[size:]
  }
}