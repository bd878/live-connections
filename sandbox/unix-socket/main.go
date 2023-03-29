package main

import (
  "net"
  "fmt"
  "time"
  "os"
  "os/signal"
  "path/filepath"
)

const socketPath = "/tmp/test.socket"

var conn net.Conn

const eof = "\n"

func produce(quit chan struct{}) {
  fmt.Println("start producer")

  go func() {
    <-quit
    conn.Close()
  }()

  for {
    n, err := conn.Write([]byte("test"))
    fmt.Printf("written %d bytes\n", n)
    if err != nil {
      fmt.Println("error while writing:", err)
    }
    time.Sleep(1000 * time.Millisecond)
  }
}

func consume(quit chan struct{}) {
  fmt.Println("start consumer")

  go func() {
    <-quit
    conn.Close()
  }()

  buf := make([]byte, 1024)

  for {
    n, err := conn.Read(buf)
    fmt.Printf("read %d bytes, received %d\n", n, buf)
    if err != nil {
      fmt.Println("failed to read")
      return
    }
  }
}

func main() {
  dir, err := os.MkdirTemp("", "unix")
  if err != nil {
    fmt.Println("failed to make tmp dir")
    return
  }
  defer func() {
    if err := os.RemoveAll(dir); err != nil {
      fmt.Println("failed to remove tmp dir")
    }
  }()

  socketPath := filepath.Join(dir, fmt.Sprintf("%d.socket", os.Getpid()))
  listener, _ := net.Listen("unix", socketPath)
  if err != nil {
    fmt.Println("failed to create listener")
    return
  }

  err = os.Chmod(socketPath, os.ModeSocket|0666)
  if err != nil {
    fmt.Println("failed to change socket permissions")
    return
  }

  go func() {
    conn, err := listener.Accept()
    if err != nil {
      return
    }

    go func() {
      for {
        buf := make([]byte, 1024)
        n, err := conn.Read(buf)
        if err != nil {
          return
        }

        _, err = conn.Write(buf[:n])
        if err != nil {
          return
        }
      }
    }()
  }()

  conn, _ = net.Dial("unix", listener.Addr().String())
  defer func() { _ = conn.Close() }()

  quit := make(chan struct{})
  defer func() { close(quit) }()

  go produce(quit)
  go consume(quit)

  go func() {
    sig := make(chan os.Signal, 1)
    signal.Notify(sig, os.Interrupt)

    <-sig

    fmt.Println("interrupt caught")
    quit <- struct{}{}
  }()

  <-quit

  fmt.Println("program quit")
}