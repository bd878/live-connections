package main

import (
  "log"

  "github.com/teralion/live-connections/disk/internal/server"
)

const addr = "localhost:50051"

func main() {
  grpcServer := server.NewGRPCServer(addr)

  log.Println("server is listening on =", addr)
  if err := grpcServer.Serve(); err != nil {
    log.Fatalf("failed to serve: %v\n", err)
  }
}