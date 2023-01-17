package main

import (
  "log"
  "os"
  "net/http"

  dotenv "github.com/joho/godotenv"

  "github.com/teralion/live-connections/server/internal/server"
)

func main() {
  if err := dotenv.Load(); err != nil {
    log.Fatal("Error loading .env file")
  }

  var serverCrt, serverKey, addr string
  var ok bool
  serverCrt, ok = os.LookupEnv("LC_SERVER_CRT_PATH")
  if !ok {
    log.Fatal("No server crt path provided")
  }

  serverKey, ok = os.LookupEnv("LC_SERVER_KEY_PATH")
  if !ok {
    log.Fatal("No client crt path provded")
  }

  addr, ok = os.LookupEnv("LC_SERVER_ADDR")
  if !ok {
    log.Fatalf("Client is lack of addr")
  }

  done := make(chan struct{})

  srv := server.NewHTTPServer(addr, done)
  log.Println("server is listening on =", addr)
  if err := srv.ListenAndServeTLS(serverCrt, serverKey); err != http.ErrServerClosed {
    log.Fatalf("HTTP server ListenAndServeTLS: %v", err)
  }

  <-done
}
