package protocol

import "time"

const (
  MaxPayloadSize int64 = 512
  PongWait = 60 * time.Second
  PingPeriod = (PongWait * 9) / 10 // a bit less than pongWait
  WriteWait = 10 * time.Second
  MaxClients = 10
  MTU = 256
)
