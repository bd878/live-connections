package server

import (
  "errors"
  "math/rand"
)

var colors = []string{
  "80FF00",
  "0000FF",
  "808080",
  "FF8000",
}

type Colors struct {
  occupied map[string]bool
}

func NewColors() *Colors {
  return &Colors{
    occupied: make(map[string]bool, len(colors)),
  }
}

func (c *Colors) HasFree() bool {
  return len(c.occupied) < len(colors)
}

func (c *Colors) GetColor() (string, error) {
  if !c.HasFree() {
    return "", errors.New("no free colors left")
  }

  color := colors[rand.Intn(len(colors))]
  for i := 0; i < len(colors) && c.IsOccupied(color); i++ {
    color = colors[rand.Intn(len(colors))]
  }
  c.occupied[color] = true

  return color, nil
}

func (c *Colors) IsOccupied(color string) bool {
  return c.occupied[color]
}

func (c *Colors) Release(color string) {
  delete(c.occupied, color)
}