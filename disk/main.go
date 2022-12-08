package main

import (
  "net"
  "time"
  "context"
  "fmt"
  "log"
  "math/rand"
  "strings"
  "os"

  "google.golang.org/grpc"
  pb "github.com/teralion/live-connections/disk/proto"
)

const (
  port = 50051
  charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
  dirLength = 10
  baseDir = filepath.Join("../files")
)

func randomString(n int) string {
  b := strings.Builder{}
  b.Grow(n)
  for i := 0; i < n; i++ {
    b.WriteByte(charset[rand.Intn(len(charset))])
  }
  return b.String()
}

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
  rand.Seed(time.Now().UnixNano())
  dirName := randomString(nameLength)

  if err := os.MkdirAll(filepath.Join(baseDir, dirName), 0750); err != nil {
    log.Printf("error creating area: %v\n", err)
    return nil, err
  }

  return &pb.CreateAreaResponse{name: dirName}, nil
}

func (s *areaManagerServer) ListUsers(ctx context.Context, request *pb.ListAreaUsersRequest) (*pb.ListAreaUsersResponse, error) {
  log.Println("Received list users request")
  return &pb.ListAreaUsersResponse{}, nil
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