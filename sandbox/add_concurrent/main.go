package main

import (
  "os"
  "runtime/debug"
  "fmt"
)

func main() {
  prev := debug.SetGCPercent(10)
  fmt.Println("prev:", prev)

  next := debug.SetGCPercent(10)
  fmt.Println("next:", next)

  key := "GOGC"
  gogc, ok := os.LookupEnv(key)
  if !ok {
    fmt.Printf("%s not set\n", key)
  } else {
    fmt.Println("now:", gogc)
  }
}
