package main

import (
  "net"
  "context"
  "path/filepath"
  "fmt"
  "log"
  "errors"
  "os"

  "google.golang.org/grpc"
  utils "github.com/teralion/live-connections/disk/utils"
  pb "github.com/teralion/live-connections/disk/proto"
)

const (
  port = 50051
  areaNameLength = 10
  userNameLength = 10
  allEntries = -1
)

var baseDir = filepath.Join("../", "files")

type areaManagerServer struct {
  pb.UnimplementedAreaManagerServer
}

type userManagerServer struct {
  pb.UnimplementedUserManagerServer
}

type shapeManagerServer struct {
  pb.UnimplementedShapeManagerServer
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

func (s *userManagerServer) Add(ctx context.Context, request *pb.AddUserRequest) (*pb.AddUserResponse, error) {
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

func main() {
  lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
  if err != nil {
    log.Fatalf("failed to listen: %v\n", err)
  }

  grpcServer := grpc.NewServer()
  pb.RegisterAreaManagerServer(grpcServer, &areaManagerServer{})
  pb.RegisterUserManagerServer(grpcServer, &userManagerServer{})
  pb.RegisterShapeManagerServer(grpcServer, &shapeManagerServer{})
  log.Printf("server listening at %v", lis.Addr())

  if err := grpcServer.Serve(lis); err != nil {
    log.Fatalf("failed to serve: %v\n", err)
  }
}