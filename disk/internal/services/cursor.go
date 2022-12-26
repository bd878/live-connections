package services

import (
  "fmt"
  "context"

  pb "github.com/teralion/live-connections/disk/pkg/proto"
)

type CursorServer struct {
  pb.UnimplementedCursorManagerServer
  Dir string
}

func NewCursorManagerServer(baseDir string) *CursorServer {
  return &CursorServer{Dir: baseDir}
}

func (s *CursorServer) Write(ctx context.Context, request *pb.WriteCursorRequest) (*pb.WriteCursorResponse, error) {
  fmt.Println("not implemented")
  return &pb.WriteCursorResponse{}, nil
}

func (s *CursorServer) Read(ctx context.Context, request *pb.ReadCursorRequest) (*pb.ReadCursorResponse, error) {
  fmt.Println("not implemented")
  return &pb.ReadCursorResponse{XPos: 0, YPos: 0}, nil
}