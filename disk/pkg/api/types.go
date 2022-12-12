package api

import (
  "net"

  "google.golang.org/grpc"
)

type GRPCServer struct {
  server *grpc.Server
  listener net.Listener
}