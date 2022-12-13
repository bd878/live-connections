package api

import (
  "net"
  "log"
  "os"
  "os/signal"
  "path/filepath"

  "google.golang.org/grpc"

  services "github.com/teralion/live-connections/disk/internal/services"
  pb "github.com/teralion/live-connections/disk/pkg/proto"
)

func NewGRPCServer(addr string) *GRPCServer {
  listener, err := net.Listen("tcp", addr)
  if err != nil {
    log.Fatalf("failed to listen: %v\n", err)
  }

  grpcServer := grpc.NewServer()
  baseDir := filepath.Join("../", "files")

  pb.RegisterAreaManagerServer(grpcServer, &services.AreaManagerServer{Dir: baseDir, NameLen: 10})
  pb.RegisterUserManagerServer(grpcServer, &services.UserManagerServer{Dir: baseDir, NameLen: 10})

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

func (s *GRPCServer) Stop() {
  s.server.Stop()
}