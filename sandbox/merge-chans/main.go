package main

import "fmt"
import "sync"

func merge(chans ...<-chan int) <-chan int {
  var wg sync.WaitGroup

  out := make(chan int)

  send := func(ch <-chan int) {
    for v := range ch {
      out <- v
    }
    wg.Done()
  }

  for _, ch := range chans {
    wg.Add(1)
    go send(ch)
  }

  go func() {
    wg.Wait()
    close(out)
  }()

  return out
}

func gen(nums ...int) <-chan int {
  out := make(chan int)
  go func() {
    for _, n := range nums {
      out <- n
    }
    close(out)
  }()
  return out
}

func main() {
  out1 := gen(3, 4)
  out2 := gen(9)
  out3 := gen(5, 6, 7, 8)

  for v := range merge(out1, out2, out3) {
    fmt.Println(v)
  }
}