package main

import (
  "os"
  "net/http"

  dotenv "github.com/joho/godotenv"

  "github.com/teralion/live-connections/server/internal/server"
  "github.com/teralion/live-connections/server/internal/meta"
)

func main() {
  if err := dotenv.Load(); err != nil {
    meta.Log().Fatal("Error loading .env file")
  }

  addr, ok := os.LookupEnv("LC_SERVER_ADDR")
  if !ok {
    meta.Log().Fatal("Client is lack of addr")
  }

  done := make(chan struct{})

  srv := server.NewHTTPServer(addr, done)
  meta.Log().Info("server is listening on =", addr)
  if err := srv.ListenAndServe(); err != http.ErrServerClosed {
    meta.Log().Error("HTTP server ListenAndServe:", err)
  }

  <-done
}
