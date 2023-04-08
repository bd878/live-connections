package main

import (
  "os"
  "fmt"

  dotenv "github.com/joho/godotenv"

  "github.com/teralion/live-connections/meta"
  "github.com/teralion/live-connections/disk/pkg/api"
)

func main() {
  if err := dotenv.Load(); err != nil {
    meta.Log().Fatal("Error loading .env file")
  }

  addr, ok := os.LookupEnv("LC_DISK_ADDR")
  if !ok {
    meta.Log().Fatal("Disk is lack of addr")
  }

  grpcServer := api.NewGRPCServer(addr)

  meta.Log().Debug("server is listening on =", addr)
  if err := grpcServer.Serve(); err != nil {
    meta.Log().Fatal(fmt.Sprintf("failed to serve: %v\n", err))
  }
}