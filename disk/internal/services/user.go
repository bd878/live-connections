package services

import (
  "context"
  "os"
  "errors"
  "path/filepath"

  utils "github.com/teralion/live-connections/disk/internal/utils"
  pb "github.com/teralion/live-connections/disk/pkg/proto"
)

type UserManagerServer struct {
  pb.UnimplementedUserManagerServer
  Dir string
  NameLen int
}

func (s *UserManagerServer) Add(ctx context.Context, request *pb.AddUserRequest) (*pb.AddUserResponse, error) {
  areaName := request.Area

  if !utils.IsNameSafe(areaName) {
    return nil, errors.New("not safe")
  }

  userName := utils.RandomString(s.NameLen)
  if err := os.Mkdir(filepath.Join(s.Dir, areaName,  userName), 0750); err != nil {
    return nil, err
  }

  return &pb.AddUserResponse{Name: userName}, nil
}
