package main

import (
  "fmt"
)

type Voice struct {
  ch chan string
}

func (c *Voice) Send() chan string {
  return c.ch
}

func main() {
  v := &Voice{
    ch: make(chan string, 1),
  }

  v.Send() <- "aaa"
  result := <-v.Send()

  fmt.Println(result)
}