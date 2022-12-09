package main

import (
  "os"
  "path/filepath"
  "log"
  "time"
  "flag"
  "fmt"
  "math/rand"
  "strings"
)

var (
  dirsCount = flag.Int("n", 100, "dirs count")
  baseDir = flag.String("dir", "base", "base dir path")
  clear = flag.Bool("rm", false, "clear base dir")
)

const (
  charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
  readdirstep = 5
)

func randomString(n int) string {
  rand.Seed(time.Now().UnixNano())

  b := strings.Builder{}
  b.Grow(n)
  for i := 0; i < n; i++ {
    b.WriteByte(charset[rand.Intn(len(charset))])
  }
  return b.String()
}

func createDir(base string) error {
  dirName := randomString(10)

  if err := os.Mkdir(filepath.Join(base, dirName), 0666); err != nil {
    return err
  }

  return nil
}

func main() {
  flag.Parse()

  if *clear {
    err := os.RemoveAll(*baseDir)
    if os.IsNotExist(err) {
      return
    } else {
      log.Fatal(err)
    }
  }

  if _, err := os.Stat(*baseDir); os.IsNotExist(err) {
    if err = os.Mkdir(*baseDir, 0750); err != nil {
      log.Fatal(err)
    }
  }

  for i := 0; i < *dirsCount; i++ {
    if err := createDir(*baseDir); err != nil {
      log.Fatal(err)
    }
  }

  files, err := os.ReadDir(*baseDir)
  if err != nil {
    log.Fatal(err)
  }

  for _, file := range files {
    fmt.Println(file.Name())
  }
}