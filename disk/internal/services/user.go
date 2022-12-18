package services

import (
  "context"
  "os"
  "log"
  "errors"
  "path/filepath"
  "bufio"

  utils "github.com/teralion/live-connections/disk/internal/utils"
  pb "github.com/teralion/live-connections/disk/pkg/proto"
)

type UserManagerServer struct {
  pb.UnimplementedUserManagerServer
  Dir string
  NameLen int
}

func NewUserManagerServer(baseDir string) *UserManagerServer {
  return &UserManagerServer{Dir: baseDir, NameLen: 10}
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

func (s *UserManagerServer) Write(ctx context.Context, request *pb.WriteUserRequest) (*pb.WriteUserResponse, error) {
  filepath.Join(s.Dir, request.Area, request.Name)

  log.Fatalf("not implemented")
  return &pb.WriteUserResponse{}, errors.New("not implemented")
}

func writeFile(fp string, p []byte) error {
  storeFile, err := os.OpenFile(
    fp,
    os.O_RDWR|os.O_CREATE|os.O_APPEND,
    0644,
  )
  if err != nil {
    return err
  }
  buf := bufio.NewWriter(storeFile)
  _, err = buf.Write(p)
  if err != nil {
    log.Printf("writeFile error write to buffer = %v\n", err)
    return err
  }
  if err := buf.Flush(); err != nil {
    log.Printf("buf.Flush error flush to file = %v\n", err)
    return err
  }
  return nil
}