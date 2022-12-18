package main

import (
  "log"
  "net/http"
  "path/filepath"

  server "github.com/teralion/live-connections/client/internal/server"
)

const addr = "localhost:8081"
var (
  serverCrt = filepath.Join("./", "cmd/server", "server.crt")
  serverKey = filepath.Join("./", "cmd/server", "server.key")
)

func main() {
  done := make(chan struct{})

  srv := server.NewHTTPServer(addr, done)
  log.Println("server is listening on =", addr)
  if err := srv.ListenAndServeTLS(serverCrt, serverKey); err != http.ErrServerClosed {
    log.Fatalf("HTTP server ListenAndServeTLS: %v", err)
  }

  <-done
}