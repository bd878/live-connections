package main

import (
  "log"

  "github.com/teralion/live-connections/disk/pkg/api"
)

const addr = "localhost:50051"

func main() {
  grpcServer := api.NewGRPCServer(addr)

  log.Println("server is listening on =", addr)
  if err := grpcServer.Serve(); err != nil {
    log.Fatalf("failed to serve: %v\n", err)
  }
}