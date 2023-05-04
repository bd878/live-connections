package main

import (
  "fmt"
  "time"
  "context"
  "sync"
)

type Val struct {
  Ch chan string
  ctx context.Context
}

func (v *Val) Loop(pr string) {
  ticker := time.NewTicker(time.Second)
  defer ticker.Stop()

  for {
    select {
    case <-v.ctx.Done():
      fmt.Println(pr, ": done")
      return
    case t := <-ticker.C:
      fmt.Println(pr, ": tick", t)
    }
  }
}

func main() {
  ctx, cancel := context.WithCancel(context.Background())

  v := &Val{
    Ch: make(chan string, 2),
    ctx: ctx,
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
  cancel()

  fmt.Println("waiting...")
  wg.Wait()
}
