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

  var serverCrt, clientCrt, addr string
  var ok bool
  serverCrt, ok = os.LookupEnv("LC_CLIENT_CRT_PATH_SERVER")
  if !ok {
    log.Fatal("No server crt path provided")
  }

  clientCrt, ok = os.LookupEnv("LC_CLIENT_CRT_PATH_CLIENT")
  if !ok {
    log.Fatal("No client crt path provded")
  }

  addr, ok = os.LookupEnv("LC_CLIENT_ADDR")
  if !ok {
    log.Fatalf("Client is lack of addr")
  }

  done := make(chan struct{})

  srv := server.NewHTTPServer(addr, done)
  log.Println("server is listening on =", addr)
  if err := srv.ListenAndServeTLS(serverCrt, clientCrt); err != http.ErrServerClosed {
    log.Fatalf("HTTP server ListenAndServeTLS: %v", err)
  }

  <-done
}