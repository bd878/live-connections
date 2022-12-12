package server

import (
  "testing"
  "time"

  "github.com/teralion/live-connections/disk/pkg/api"
)

func TestGRPCConnection(t *testing.T) {
  grpcServer := api.NewGRPCServer("localhost:50051")
  go grpcServer.Serve()
  defer grpcServer.Stop()
  t.Log("server is running")

  time.Sleep(2 * time.Second)
}
