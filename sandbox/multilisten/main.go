package main

import (
  "fmt"
  "time"
  "sync"
)

type Val struct {
  Ch chan string
}

func (v *Val) Loop(pr string) {
  ticker := time.NewTicker(time.Second)
  defer ticker.Stop()

  for {
    select {
    case str := <-v.Ch:
      fmt.Println(pr, ": ", str)
      return
    case t := <-ticker.C:
      fmt.Println(pr, ": tick", t)
    }
  }
}

func main() {
  v := &Val{
    Ch: make(chan string, 2),
  }

  var wg sync.WaitGroup
  go func() {
    defer wg.Done()
    wg.Add(1)
    v.Loop("first")
  }()
  go func() {
    defer wg.Done()
    wg.Add(1)
    v.Loop("second")
  }()

  fmt.Println("sleeping...")
  time.Sleep(3 * time.Second)

  v.Ch <- "aaa test"
  v.Ch <- "bbb test"

  fmt.Println("waiting...")
  wg.Wait()
}
