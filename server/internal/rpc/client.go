package rpc

import (
  "net/http"
  "io"
  "os"
  "time"
  "fmt"
  "context"
  "strings"

  "google.golang.org/grpc"
  "github.com/gorilla/mux"

  "github.com/teralion/live-connections/server/proto/disk"
  "github.com/teralion/live-connections/meta"
)

type Disk struct {
  timeout time.Duration
  area disk.AreaManagerClient
  user disk.UserManagerClient
  texts disk.TextsManagerClient
  square disk.SquareManagerClient
}

type Coords struct {
  XPos float32
  YPos float32
}

func NewDisk() *Disk {
  opts := []grpc.DialOption{grpc.WithInsecure()}

  serverAddr, ok := os.LookupEnv("LC_DISK_ADDR")
  if !ok {
    meta.Log().Fatal("Disk is lack of addr")
  }

  timeoutStr, ok := os.LookupEnv("LC_DISK_REQUEST_TIMEOUT")
  if !ok {
    meta.Log().Fatal("No disk request timeout specified")
  }

  diskRequestTimeout, err := time.ParseDuration(timeoutStr)
  if  err != nil {
    meta.Log().Fatal("Failed to parse timeout duration")
  }

  conn, err := grpc.Dial(serverAddr, opts...)
  if err != nil {
    meta.Log().Fatal("failed to dial: %v\n", err)
  }

  area := disk.NewAreaManagerClient(conn)
  user := disk.NewUserManagerClient(conn)
  texts := disk.NewTextsManagerClient(conn)
  square := disk.NewSquareManagerClient(conn)

  return &Disk{
    timeout: diskRequestTimeout,
    area: area,
    user: user,
    texts: texts,
    square: square,
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
    meta.Log().Fatal("user.Add failed: %v", err)
  }
  w.Header().Set("Content-Type", "text/plain; charset=utf-8")
  w.Header().Set("Access-Control-Allow-Origin", "*")
  fmt.Fprint(w, resp.Name)
}

func (d *Disk) HandleNewArea(w http.ResponseWriter, r *http.Request) {
  meta.Log().Debug("new area request")

  ctx, cancel := context.WithTimeout(context.Background(), d.timeout)
  defer cancel()
  resp, err := d.area.Create(ctx, &disk.CreateAreaRequest{})
  if err != nil {
    meta.Log().Fatal("area.Create failed: %v", err)
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
    meta.Log().Fatal("area.ListUsers failed: %v", err)
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
    meta.Log().Fatal("area.HasUser failed: %v", err)
  }
  return resp.Result
}

func (d *Disk) WriteSquareCoords(area, user string, XPos, YPos float32) {
  ctx, cancel := context.WithTimeout(context.Background(), d.timeout)
  defer cancel()

  coords := &disk.Coords{XPos: XPos, YPos: YPos}
  _, err := d.square.Write(ctx, &disk.WriteSquareRequest{Area: area, Name: user, Coords: coords})
  if err != nil {
    meta.Log().Fatal("square.Write failed: %v", err)
  }
}

func (d *Disk) WriteText(area, user, value string) {
  ctx, cancel := context.WithTimeout(context.Background(), d.timeout)
  defer cancel()

  text := &disk.Text{Value: value}
  _, err := d.texts.Write(ctx, &disk.WriteTextRequest{Area: area, Name: user, Text: text})
  if err != nil {
    meta.Log().Fatal("text.Write failed: %v", err)
  }
}

func (d *Disk) ReadSquareCoords(area, user string) (*Coords, error) {
  ctx, cancel := context.WithTimeout(context.Background(), d.timeout)
  defer cancel()

  resp, err := d.square.Read(ctx, &disk.ReadRequest{Area: area, Name: user})
  if err != nil {
    meta.Log().Warn("square.Read failed: %v", err)
    return nil, err
  }

  return &Coords{resp.XPos, resp.YPos}, nil
}

func (d *Disk) ReadText(area, user string) (string, error) {
  ctx, cancel := context.WithTimeout(context.Background(), d.timeout)
  defer cancel()

  resp, err := d.texts.Read(ctx, &disk.ReadRequest{Area: area, Name: user})
  if err != nil {
    meta.Log().Warn("text.Read failed: %v", err)
    return "", err
  }

  return resp.Value, nil
}
