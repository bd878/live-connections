package rpc

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
  "github.com/teralion/live-connections/server/internal/types"
)

type Disk struct {
  timeout time.Duration
  area disk.AreaManagerClient
  user disk.UserManagerClient
  cursor disk.CursorManagerClient
}

func NewDisk() *Disk {
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

  area := disk.NewAreaManagerClient(conn)
  user := disk.NewUserManagerClient(conn)
  cursor := disk.NewCursorManagerClient(conn)

  return &Disk{
    timeout: diskRequestTimeout,
    area: area,
    user: user,
    cursor: cursor,
  }
}

func (d *Disk) HandleJoin(w http.ResponseWriter, r *http.Request) {
  ctx, cancel := context.WithTimeout(context.Background(), d.timeout)
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

  resp, err := d.user.Add(ctx, &disk.AddUserRequest{Area: areaName.String()})
  if err != nil {
    log.Fatalf("user.Add failed: %v", err)
  }
  w.Header().Set("Content-Type", "text/plain; charset=utf-8")
  w.Header().Set("Access-Control-Allow-Origin", "*")
  fmt.Fprint(w, resp.Name)
}

func (d *Disk) HandleNewArea(w http.ResponseWriter, r *http.Request) {
  log.Println("new area request")

  ctx, cancel := context.WithTimeout(context.Background(), d.timeout)
  defer cancel()
  resp, err := d.area.Create(ctx, &disk.CreateAreaRequest{})
  if err != nil {
    log.Fatalf("area.Create failed: %v", err)
  }
  w.Header().Set("Content-Type", "text/plain; charset=utf-8")
  w.Header().Set("Access-Control-Allow-Origin", "*")
  fmt.Fprint(w, resp.Name)
}

func (d *Disk) HandleAreaUsers(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  p := vars["id"]

  ctx, cancel := context.WithTimeout(context.Background(), d.timeout)
  defer cancel()
  resp, err := d.area.ListUsers(ctx, &disk.ListAreaUsersRequest{Name: p})
  if err != nil {
    log.Fatalf("area.ListUsers failed: %v", err)
  }
  w.Header().Set("Access-Control-Allow-Origin", "*")
  w.Header().Set("Content-Type", "text/plain; charset=utf-8")
  fmt.Fprint(w, resp.GetUsers())
}

func (d *Disk) HasUser(area string, user string) bool {
  ctx, cancel := context.WithTimeout(context.Background(), d.timeout)
  defer cancel()
  resp, err := d.area.HasUser(ctx, &disk.HasUserRequest{Area: area, User: user})
  if err != nil {
    log.Fatalf("area.HasUser failed: %v", err)
  }
  return resp.Result
}

func (d *Disk) WriteMouseCoords(area string, user string, xPos float32, yPos float32) {
  ctx, cancel := context.WithTimeout(context.Background(), d.timeout)
  defer cancel()

  coords := &disk.Coords{XPos: xPos, YPos: yPos}
  _, err := d.cursor.Write(ctx, &disk.WriteCursorRequest{Area: area, Name: user, Coords: coords})
  if err != nil {
    log.Fatalf("cursor.Write failed: %v", err)
  }
}

func (d *Disk) ReadMouseCoords(area string, user string) *types.Coords {
  ctx, cancel := context.WithTimeout(context.Background(), d.timeout)
  defer cancel()

  resp, err := d.cursor.Read(ctx, &disk.ReadCursorRequest{Area: area, Name: user})
  if err != nil {
    log.Println("cursor.Read failed: %v", err)
    return nil
  }

  coords := types.Coords{XPos: resp.XPos, YPos: resp.YPos}
  return &coords
}