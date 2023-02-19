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

  var serverCrt, serverKey, addr string
  var ok bool
  serverCrt, ok = os.LookupEnv("LC_SERVER_CRT_PATH")
  if !ok {
    meta.Log().Fatal("No server crt path provided")
  }

  serverKey, ok = os.LookupEnv("LC_SERVER_KEY_PATH")
  if !ok {
    meta.Log().Fatal("No client crt path provded")
  }

  addr, ok = os.LookupEnv("LC_SERVER_ADDR")
  if !ok {
    meta.Log().Fatal("Client is lack of addr")
  }

  done := make(chan struct{})

  srv := server.NewHTTPServer(addr, done)
  meta.Log().Info("server is listening on =", addr)
  if err := srv.ListenAndServeTLS(serverCrt, serverKey); err != http.ErrServerClosed {
    meta.Log().Error("HTTP server ListenAndServeTLS:", err)
  }

  <-done
}
