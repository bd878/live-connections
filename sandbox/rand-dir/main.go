package main

import (
  "math/rand"
  "time"
  "os"
  "strings"
  "log"
  "path/filepath"
  "flag"
)

const (
  charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

var (
  clear = flag.Bool("clear", false, "empty base dir and remove it too")
  nameLength = flag.Int("n", 10, "dir name length")
  baseDir = flag.String("dir", "./dirs", "base dir to nest new dirs")
)

func randomString(n int) string {
  b := strings.Builder{}
  b.Grow(n)
  for i := 0; i < n; i++ {
    b.WriteByte(charset[rand.Intn(len(charset))])
  }
  return b.String()
}

func main() {
  flag.Parse()

  if *clear && *baseDir != "" {
    if err := os.RemoveAll(*baseDir); err != nil {
      log.Fatal(err)
    } else {
      log.Println("ok")
    }
  } else {
    rand.Seed(time.Now().UnixNano())
    dirName := randomString(*nameLength)

    if err := os.MkdirAll(filepath.Join(*baseDir, dirName), 0750); err != nil {
      log.Fatal(err)
    } else {
      log.Println("ok")
    }
  }
}