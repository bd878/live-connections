package main

import (
  "log"
  "os"
  "net/http"

  dotenv "github.com/joho/godotenv"

  server "github.com/teralion/live-connections/client/internal/server"
)

func main() {
  if err := dotenv.Load(); err != nil {
    log.Fatal("Error loading .env file")
  }

  addr, ok := os.LookupEnv("LC_CLIENT_ADDR")
  if !ok {
    log.Fatalf("Client is lack of addr")
  }

  done := make(chan struct{})

  srv := server.NewHTTPServer(addr, done)
  log.Println("server is listening on =", addr)
  if err := srv.ListenAndServe(); err != http.ErrServerClosed {
    log.Fatalf("HTTP server ListenAndServe: %v", err)
  }

  <-done
}