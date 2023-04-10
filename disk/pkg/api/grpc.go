package api

import (
  "net"
  "log"
  "os"
  "os/signal"
  "path/filepath"

  "google.golang.org/grpc"

  services "github.com/bd878/live-connections/disk/internal/services"
  pb "github.com/bd878/live-connections/disk/pkg/proto"
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
  baseDir := filepath.Join("../", "files")

  pb.RegisterAreaManagerServer(grpcServer, services.NewAreaManagerServer(baseDir))
  pb.RegisterUserManagerServer(grpcServer, services.NewUserManagerServer(baseDir))
  pb.RegisterSquareManagerServer(grpcServer, services.NewSquareManagerServer(baseDir))
  pb.RegisterTextsManagerServer(grpcServer, services.NewTextsManagerServer(baseDir))

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