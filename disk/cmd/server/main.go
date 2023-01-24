package main

import (
  "log"
  "os"

  dotenv "github.com/joho/godotenv"

  "github.com/teralion/live-connections/disk/pkg/api"
)

func main() {
  if err := dotenv.Load(); err != nil {
    log.Fatal("Error loading .env file")
  }

  var addr string
  var ok bool
  addr, ok = os.LookupEnv("LC_DISK_ADDR")
  if !ok {
    log.Fatalf("Disk is lack of addr")
  }

  grpcServer := api.NewGRPCServer(addr)

  log.Println("server is listening on =", addr)
  if err := grpcServer.Serve(); err != nil {
    log.Fatalf("failed to serve: %v\n", err)
  }
}