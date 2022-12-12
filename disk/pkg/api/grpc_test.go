package api

import (
  "testing"
)

func TestGRPCServer(t *testing.T) {
  grpcServer := NewGRPCServer("localhost:50051")
  grpcServer.Stop()
}