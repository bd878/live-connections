package conn

import (
  "net/http"
  "io"
  "os"
  "time"
  "log"
  "fmt"
  "context"
  "strings"

  "google.golang.org/grpc"
  "github.com/gorilla/mux"

  "github.com/teralion/live-connections/server/proto/disk"
  t "github.com/teralion/live-connections/server/internal/types"
)

type LiveConnections struct {
  diskTimeout time.Duration
  areaClient disk.AreaManagerClient
  userClient disk.UserManagerClient
  cursorClient disk.CursorManagerClient
}

func NewLiveConnections() *LiveConnections {
  opts := []grpc.DialOption{grpc.WithInsecure()}

  serverAddr, ok := os.LookupEnv("LC_DISK_ADDR")
  if !ok {
    log.Fatalf("Disk is lack of addr")
  }

  timeoutStr, ok := os.LookupEnv("LC_DISK_REQUEST_TIMEOUT")
  if !ok {
    log.Fatalf("No disk request timeout specified")
  }

  diskRequestTimeout, err := time.ParseDuration(timeoutStr)
  if  err != nil {
    log.Fatalf("Failed to parse timeout duration")
  }

  conn, err := grpc.Dial(serverAddr, opts...)
  if err != nil {
    log.Fatalf("failed to dial: %v\n", err)
  }

  areaClient := disk.NewAreaManagerClient(conn)
  userClient := disk.NewUserManagerClient(conn)
  cursorClient := disk.NewCursorManagerClient(conn)

  return &LiveConnections{
    diskTimeout: diskRequestTimeout,
    areaClient: areaClient,
    userClient: userClient,
    cursorClient: cursorClient,
  }
}

func (s *LiveConnections) HandleJoin(w http.ResponseWriter, r *http.Request) {
  ctx, cancel := context.WithTimeout(context.Background(), s.diskTimeout)
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

  resp, err := s.userClient.Add(ctx, &disk.AddUserRequest{Area: areaName.String()})
  if err != nil {
    log.Fatalf("userClient.Add failed: %v", err)
  }
  w.Header().Set("Content-Type", "text/plain; charset=utf-8")
  w.Header().Set("Access-Control-Allow-Origin", "*")
  fmt.Fprint(w, resp.Name)
}

func (s *LiveConnections) HandleNewArea(w http.ResponseWriter, r *http.Request) {
  log.Println("new area request")

  ctx, cancel := context.WithTimeout(context.Background(), s.diskTimeout)
  defer cancel()
  resp, err := s.areaClient.Create(ctx, &disk.CreateAreaRequest{})
  if err != nil {
    log.Fatalf("areaClient.Create failed: %v", err)
  }
  w.Header().Set("Content-Type", "text/plain; charset=utf-8")
  w.Header().Set("Access-Control-Allow-Origin", "*")
  fmt.Fprint(w, resp.Name)
}

func (s *LiveConnections) HandleAreaUsers(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  p := vars["id"]

  ctx, cancel := context.WithTimeout(context.Background(), s.diskTimeout)
  defer cancel()
  resp, err := s.areaClient.ListUsers(ctx, &disk.ListAreaUsersRequest{Name: p})
  if err != nil {
    log.Fatalf("areaClient.ListUsers failed: %v", err)
  }
  w.Header().Set("Access-Control-Allow-Origin", "*")
  w.Header().Set("Content-Type", "text/plain; charset=utf-8")
  fmt.Fprint(w, resp.GetUsers())
}

func (s *LiveConnections) HasUser(area string, user string) bool {
  ctx, cancel := context.WithTimeout(context.Background(), s.diskTimeout)
  defer cancel()
  resp, err := s.areaClient.HasUser(ctx, &disk.HasUserRequest{Area: area, User: user})
  if err != nil {
    log.Fatalf("areaClient.HasUser failed: %v", err)
  }
  return resp.Result
}

func (s *LiveConnections) WriteMouseCoords(area string, user string, xPos float32, yPos float32) {
  ctx, cancel := context.WithTimeout(context.Background(), s.diskTimeout)
  defer cancel()

  coords := &disk.Coords{XPos: xPos, YPos: yPos}
  _, err := s.cursorClient.Write(ctx, &disk.WriteCursorRequest{Area: area, Name: user, Coords: coords})
  if err != nil {
    log.Fatalf("cursorClient.Write failed: %v", err)
  }
}

func (s *LiveConnections) ReadMouseCoords(area string, user string) *t.Coords {
  ctx, cancel := context.WithTimeout(context.Background(), s.diskTimeout)
  defer cancel()

  resp, err := s.cursorClient.Read(ctx, &disk.ReadCursorRequest{Area: area, Name: user})
  if err != nil {
    log.Println("cursorClient.Read failed: %v", err)
    return nil
  }

  coords := t.Coords{XPos: resp.XPos, YPos: resp.YPos}
  return &coords
}