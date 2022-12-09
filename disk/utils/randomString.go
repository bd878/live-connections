package utils

import (
  "time"
  "math/rand"
  "strings"
)

const (
  charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

func RandomString(n int) string {
  rand.Seed(time.Now().UnixNano())

  b := strings.Builder{}
  b.Grow(n)
  for i := 0; i < n; i++ {
    b.WriteByte(charset[rand.Intn(len(charset))])
  }
  return b.String()
}
