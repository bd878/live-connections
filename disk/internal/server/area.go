package server

import (
  "context"
  "os"
  "log"
  "path/filepath"
  "errors"

  utils "github.com/teralion/live-connections/disk/internal/utils"
  pb "github.com/teralion/live-connections/disk/api/v1"
)

const (
  areaNameLength = 10
  allEntries = -1
)

type areaManagerServer struct {
  pb.UnimplementedAreaManagerServer
}

func (s *areaManagerServer) Create(ctx context.Context, request *pb.CreateAreaRequest) (*pb.CreateAreaResponse, error) {
  areaName := utils.RandomString(areaNameLength)

  if !utils.IsNameSafe(areaName) {
    log.Printf("not safe name", areaName)
    return nil, errors.New("not safe")
  }

  if err := os.MkdirAll(filepath.Join(baseDir, areaName), 0750); err != nil {
    log.Printf("error creating area: %v\n", err)
    return nil, err
  }

  return &pb.CreateAreaResponse{Name: areaName}, nil
}

func (s *areaManagerServer) ListUsers(ctx context.Context, request *pb.ListAreaUsersRequest) (*pb.ListAreaUsersResponse, error) {
  if !utils.IsNameSafe(request.Name) {
    log.Printf("not safe name", request.Name)
    return nil, errors.New("not safe")
  }

  dir, err := os.Open(filepath.Join(baseDir, request.Name))
  if err != nil {
    if os.IsNotExist(err) {
      return &pb.ListAreaUsersResponse{Users: []string{}}, os.ErrNotExist
    } else {
      return &pb.ListAreaUsersResponse{Users: []string{}}, errors.New("read dir error")
    }
  }

  names, err := dir.Readdirnames(allEntries)
  if err != nil {
    return &pb.ListAreaUsersResponse{Users: names}, errors.New("read entries error")
  }

  log.Printf("% v\n", names)
  return &pb.ListAreaUsersResponse{Users: names}, nil
}
