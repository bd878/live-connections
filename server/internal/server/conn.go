package server

import (
  "net/http"
  "io"
  "time"
  "log"
  "fmt"
  "context"
  "strings"

  "google.golang.org/grpc"

  "github.com/teralion/live-connections/server/proto/disk"
)

const (
  diskRequestTimeout = 10*time.Second
  serverAddr = "localhost:50051"
)

type liveConnections struct {
  AreaClient disk.AreaManagerClient
  UserClient disk.UserManagerClient
}

func NewLiveConnections() *liveConnections {
  opts := []grpc.DialOption{grpc.WithInsecure()}

  conn, err := grpc.Dial(serverAddr, opts...)
  if err != nil {
    log.Fatalf("failed to dial: %v\n", err)
  }
  defer conn.Close()

  areaClient := disk.NewAreaManagerClient(conn)
  userClient := disk.NewUserManagerClient(conn)

  return &liveConnections{
    AreaClient: areaClient,
    UserClient: userClient,
  }
}

func (s *liveConnections) handleJoin(w http.ResponseWriter, r *http.Request) {
  ctx, cancel := context.WithTimeout(context.Background(), diskRequestTimeout)
  defer cancel()

  body, err := io.ReadAll(r.Body)
  if err != nil {
    http.Error(w,
      fmt.Sprint("cannot read body"),
      http.StatusBadRequest,
    )
  }
  var areaName strings.Builder
  areaName.Write(body)

  resp, err := s.UserClient.Add(ctx, &disk.AddUserRequest{Area: areaName.String()})
  if err != nil {
    log.Fatalf("userClient.Add failed: %v", err)
  }
  fmt.Fprintf(w, "%v\n", resp.Name)
}

func (s *liveConnections) handleNewArea(w http.ResponseWriter, r *http.Request) {
  ctx, cancel := context.WithTimeout(context.Background(), diskRequestTimeout)
  defer cancel()
  resp, err := s.AreaClient.Create(ctx, &disk.CreateAreaRequest{})
  if err != nil {
    log.Fatalf("areaClient.Create failed: %v", err)
  }
  fmt.Fprintln(w, resp.Name)
}

func (s *liveConnections) handleAreaUsers(w http.ResponseWriter, r *http.Request) {
  p := r.URL.Path
  fmt.Printf("area path: %v\n", p)

  ctx, cancel := context.WithTimeout(context.Background(), diskRequestTimeout)
  defer cancel()
  resp, err := s.AreaClient.ListUsers(ctx, &disk.ListAreaUsersRequest{Name: p})
  if err != nil {
    log.Fatalf("areaClient.ListUsers failed: %v", err)
  }
  fmt.Printf("% v\n", resp.GetUsers())
  fmt.Fprintf(w, "%v\n", resp.GetUsers())
}
