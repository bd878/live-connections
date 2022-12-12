package server

import (
  "testing"
  "time"

  api "github.com/teralion/live-connections/disk/pkg/api"
)

func TestServer(t *testing.T) {
  server := api.NewGRPCServer("localhost:50051")
  go server.Serve()
  defer server.Stop()

  t.Log("serving...\n")
  time.Sleep(2*time.Second)
  t.Log("stopped")
}