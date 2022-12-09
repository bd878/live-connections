package main

import (
   "io"
  "log"
  "fmt"
  "os"
  "bufio"
)

func main() {
  file, _ := os.OpenFile("./test.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)
  defer file.Close()

  info, _ := os.Stat(file.Name())
  log.Println(info.Size())

  writer2 := bufio.NewWriter(os.Stdout)

  p, _ := io.ReadAll(file)
  fmt.Printf("%s", p)

  writer2.Write(p)
  writer2.Flush()
}
