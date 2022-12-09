package main

import (
  "time"
  "fmt"
  "math/rand"
  "path/filepath"
  "os"
  "log"
  "flag"
  "bufio"
)

var (
  clear = flag.Bool("rm", false, "clear files")
  baseDir = flag.String("dir", "./files", "base files dir")
)

func getNextFilename() int64 {
  return time.Now().UnixNano()
}

type UserFile struct {
  *os.File
  buf *bufio.Writer
}

func main() {
  flag.Parse()

  if *clear && *baseDir != "" {
    if err := os.RemoveAll(*baseDir); err != nil {
      log.Fatal(err)
      return
    }
    log.Println("ok")
    return
  }

  if _, err := os.Stat(*baseDir); err != nil {
    if os.IsNotExist(err) {
      if dirErr := os.Mkdir(*baseDir, 0750); dirErr != nil {
        log.Fatal(dirErr)
      }
    } else {
      log.Fatal(err)
    }
  }

  filename := getNextFilename()
  file, err := os.OpenFile(
    filepath.Join(*baseDir, fmt.Sprintf("%d", filename)),
    os.O_RDWR|os.O_CREATE|os.O_APPEND,
    0644,
  )
  if err != nil {
    log.Fatal(err)
  }

  f := &UserFile{File: file, buf: bufio.NewWriter(file)}
  payload := make([]byte, 1<<8) // 1 byte
  rand.Seed(time.Now().UnixNano())
  if _, err := rand.Read(payload); err != nil {
    log.Fatal(err)
  }

  log.Printf("bytes: %d\n", len(payload))
  if _, err := f.buf.Write(payload); err != nil {
    log.Fatal(err)
  }

  if err := f.buf.Flush(); err != nil {
    log.Fatal(err)
  }

  log.Println("ok")
}