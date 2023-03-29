package main

import "fmt"

func producer(nums ...int) <-chan int {
  out := make(chan int)
  go func() {
    for _, n := range nums {
      out <- n
    }
    close(out)
  }()
  return out
}

func mul(in <-chan int) <-chan int {
  out := make(chan int)
  go func() {
    for n1 := range in {
      n2 := <- in
      out <- n1 * n2
    }
    close(out)
  }()
  return out
}

func main() {
  gen := producer(2, 3)
  out := mul(gen)

  fmt.Println(<-out)
}