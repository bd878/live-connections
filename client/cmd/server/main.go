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

  switch env := os.Getenv("LC_CLIENT_ENV"); env {
  case "development":
    runDev()
  case "production":
    runProd()
  default:
    log.Fatalf("unknown env:", env)
  }
}

func runDev() {
  var serverCrt, serverKey, addr string
  var ok bool
  serverCrt, ok = os.LookupEnv("LC_CLIENT_CRT_PATH")
  if !ok {
    log.Fatal("No server crt path provided")
  }

  serverKey, ok = os.LookupEnv("LC_CLIENT_KEY_PATH")
  if !ok {
    log.Fatal("No server key path provided")
  }

  addr, ok = os.LookupEnv("LC_CLIENT_ADDR")
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

func runProd() {
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