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
  dirLength = 10
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
  dirName := utils.RandomString(dirLength)

  if err := os.MkdirAll(filepath.Join(baseDir, dirName), 0750); err != nil {
    log.Printf("error creating area: %v\n", err)
    return nil, err
  }

  return &pb.CreateAreaResponse{Name: dirName}, nil
}

func (s *areaManagerServer) ListUsers(ctx context.Context, request *pb.ListAreaUsersRequest) (*pb.ListAreaUsersResponse, error) {
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

  return &pb.ListAreaUsersResponse{Users: names}, nil
}

func (s *areaManagerServer) Destroy(ctx context.Context, request *pb.DestroyAreaRequest) (*pb.DestroyAreaResponse, error) {
  log.Println("Received destroy request")
  return &pb.DestroyAreaResponse{}, nil
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