package services

import (
  "context"
  "log"
  "os"
  "errors"
  "path/filepath"

  utils "github.com/teralion/live-connections/disk/internal/utils"
  pb "github.com/teralion/live-connections/disk/api/v1"
)

const userNameLength = 10

type UserManagerServer struct {
  pb.UnimplementedUserManagerServer
}

func (s *UserManagerServer) Add(ctx context.Context, request *pb.AddUserRequest) (*pb.AddUserResponse, error) {
  areaName := request.Area

  if !utils.IsNameSafe(areaName) {
    log.Printf("not safe name", areaName)
    return nil, errors.New("not safe")
  }

  userName := utils.RandomString(userNameLength)
  if err := os.Mkdir(filepath.Join(baseDir, areaName,  userName), 0750); err != nil {
    log.Printf("error creating user: %v\n", err)
    return nil, err
  }

  return &pb.AddUserResponse{Name: userName}, nil
}
