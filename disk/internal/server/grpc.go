package server

import (
  "net"
  "log"
  "os"
  "os/signal"

  "google.golang.org/grpc"
  pb "github.com/teralion/live-connections/disk/api/v1"
)

type GRPCServer struct {
  server *grpc.Server
  listener net.Listener
}

func NewGRPCServer(addr string) *GRPCServer {
  listener, err := net.Listen("tcp", addr)
  if err != nil {
    log.Fatalf("failed to listen: %v\n", err)
  }

  grpcServer := grpc.NewServer()

  pb.RegisterAreaManagerServer(grpcServer, &areaManagerServer{})
  pb.RegisterUserManagerServer(grpcServer, &userManagerServer{})

  return &GRPCServer{listener: listener, server: grpcServer}
}

func (s *GRPCServer) Serve() error {
  done := make(chan struct{})

  go func() {
    sigint := make(chan os.Signal, 1)
    signal.Notify(sigint, os.Interrupt)
    <-sigint

    log.Println("SIGING caught")
    s.server.Stop()
    close(done)
  }()

  err := s.server.Serve(s.listener)

  <-done

  return err
}