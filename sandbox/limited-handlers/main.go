package main

import (
  "fmt"
  "time"
)

const handlersCount = 1
const maxConns = 10
const requestsNum = 100

var queue = make(chan int, maxConns)

func main() {
  serve(queue)
  producer()

  fmt.Println("main exited")
}

func producer() {
  fmt.Println("producing")

  for i := 0; i < requestsNum; i++ {
    queue <- i
  }
}

func serve(q chan int) {
  fmt.Println("creating goroutines...")

  for i := 0; i < handlersCount; i++ {
    go handle(i, q)
  }
}

func handle(num int, q chan int) {
  for v := range q {
    process(num, v)
  }
}

func process(num int, v int) {
  d, _ := time.ParseDuration("1s")
  time.Sleep(d)
  fmt.Printf("#%d: %d\n", num, v)
}