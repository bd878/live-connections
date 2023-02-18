package server

import (
  "net/http"
  "time"
  "strings"

  ws "github.com/gorilla/websocket"
)

func checkClientOrigin(r *http.Request) bool {
  origin := r.Host
  return strings.Contains(origin, "localhost")
}

var upgrader = ws.Upgrader{
  HandshakeTimeout: 10 * time.Second,
  ReadBufferSize: 512,
  WriteBufferSize: 512,
  CheckOrigin: checkClientOrigin,
}

func UpgradeConnection(w http.ResponseWriter, r *http.Request) (*ws.Conn, error) {
  conn, err := upgrader.Upgrade(w, r, nil)
  if err != nil {
    w.WriteHeader(http.StatusBadRequest)
    return nil, err
  }

  return conn, nil
}
