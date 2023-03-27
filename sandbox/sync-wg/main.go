package main

import (
  "fmt"
  "sync"
  "time"
)

func main() {
  var wg sync.WaitGroup
  wg.Add(1)
  go func() {
    defer wg.Done()
    d, _ := time.ParseDuration("2s")
    fmt.Println("blocked")
    time.Sleep(d)
    fmt.Println("ready")

  }()
  wg.Wait()
  return
}