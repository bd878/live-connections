package main

import (
  "os"
  "strings"
  "os/exec"
  "fmt"
)

func main() {
  cmd := exec.Command("ip", "a")
  var out strings.Builder
  cmd.Stdout = &out
  err := cmd.Run()
  if err != nil {
    fmt.Println(err)
    os.Exit(1)
  }
  fmt.Printf("result:\n%s\n", out.String())
}