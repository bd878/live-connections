package main

import (
  "os"
  "log"
  "io"
  "bufio"
)

func readFile() {

}

func writeFile() {

}

func main() {
  fw, err := os.OpenFile("./main.go", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0770)
  if err != nil {
    log.Fatal(err)
  }

  fr, err := os.OpenFile("./main.go", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0770)
  if err != nil {
    log.Fatal(err)
  }

  writer := bufio.NewWriter(fw)
  writer.Write([]byte("test"))
  writer.Flush()

  n, err := io.Copy(os.Stdout, fr)
  if err != nil {
    log.Fatal(err)
  }

  log.Println("copied =", n)
}