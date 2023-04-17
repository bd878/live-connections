package main

import (
  "log"
  "fmt"
  "os"
  "path/filepath"
  "os/exec"
  "strings"
)

func main() {
  base, err := os.Getwd()
  if err != nil {
    log.Fatal(err)
  }

  pSQL := filepath.Join(base, "pdns.sqlite3.sql")
  pDB := filepath.Join(base, "pdns.sqlite3")

  cmd := exec.Command("/usr/bin/sqlite3", pDB)
  fmt.Println(cmd.Args)

  in := strings.NewReader(fmt.Sprintf(".read %s", pSQL))
  cmd.Stdin = in

  err = cmd.Start()
  if err != nil {
    log.Fatal(err)
  }

  err = cmd.Wait()
  if err != nil {
    log.Fatal(err)
  }

  fmt.Println("success")
}