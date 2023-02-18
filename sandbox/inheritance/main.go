package main

import "log"

type Message struct {
  text string
}

type AuthMessage struct {
  *Message
}

func main() {
  authMsg := AuthMessage{}
  authMsg.Message = &Message{"test"}

  log.Println(authMsg.Message.text)
}