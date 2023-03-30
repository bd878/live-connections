package rpc

import (
  "time"
  "fmt"
  "context"

  "google.golang.org/grpc"

  "github.com/teralion/live-connections/disk/pkg/proto"

  "github.com/teralion/live-connections/meta"
  "github.com/teralion/live-connections/server/pkg/messages"
)

type Disk struct {
  timeout time.Duration
  area proto.AreaManagerClient
  user proto.UserManagerClient
  texts proto.TextsManagerClient
  square proto.SquareManagerClient
}

func NewDisk(addr string) *Disk {
  opts := []grpc.DialOption{grpc.WithInsecure()}

  diskRequestTimeout, err := time.ParseDuration("10s")
  if  err != nil {
    meta.Log().Fatal("Failed to parse timeout duration")
  }

  conn, err := grpc.Dial(addr, opts...)
  if err != nil {
    meta.Log().Fatal(fmt.Sprintf("failed to dial: %v\n", err))
  }

  area := proto.NewAreaManagerClient(conn)
  user := proto.NewUserManagerClient(conn)
  texts := proto.NewTextsManagerClient(conn)
  square := proto.NewSquareManagerClient(conn)

  return &Disk{
    timeout: diskRequestTimeout,
    area: area,
    user: user,
    texts: texts,
    square: square,
  }
}

func (d *Disk) CreateNewUser(area string) string {
  ctx, cancel := context.WithTimeout(context.Background(), d.timeout)
  defer cancel()

  resp, err := d.user.Add(ctx, &proto.AddUserRequest{Area: area})
  if err != nil {
    meta.Log().Fatal(fmt.Sprintf("user.Add failed: %v", err))
  }

  return resp.Name
}

func (d *Disk) CreateNewArea() string {
  ctx, cancel := context.WithTimeout(context.Background(), d.timeout)

  defer cancel()
  resp, err := d.area.Create(ctx, &proto.CreateAreaRequest{})
  if err != nil {
    meta.Log().Fatal(fmt.Sprintf("area.Create failed: %v", err))
  }

  return resp.Name
}

func (d *Disk) ListUsers(area string) []string {
  ctx, cancel := context.WithTimeout(context.Background(), d.timeout)
  defer cancel()

  resp, err := d.area.ListUsers(ctx, &proto.ListAreaUsersRequest{Name: area})
  if err != nil {
    meta.Log().Fatal(fmt.Sprintf("area.ListUsers failed: %v", err))
  }

  return resp.GetUsers()
}

func (d *Disk) HasUser(area string, user string) bool {
  ctx, cancel := context.WithTimeout(context.Background(), d.timeout)
  defer cancel()
  resp, err := d.area.HasUser(ctx, &proto.HasUserRequest{Area: area, User: user})
  if err != nil {
    meta.Log().Fatal(fmt.Sprintf("area.HasUser failed: %v", err))
  }
  return resp.Result
}

func (d *Disk) WriteSquareCoords(area, user string, XPos, YPos float32) {
  ctx, cancel := context.WithTimeout(context.Background(), d.timeout)
  defer cancel()

  coords := &proto.Coords{XPos: XPos, YPos: YPos}
  _, err := d.square.Write(ctx, &proto.WriteSquareRequest{Area: area, Name: user, Coords: coords})
  if err != nil {
    meta.Log().Fatal(fmt.Sprintf("square.Write failed: %v", err))
  }
}

func (d *Disk) WriteText(area, user, value string) {
  ctx, cancel := context.WithTimeout(context.Background(), d.timeout)
  defer cancel()

  text := &proto.Text{Value: value}
  _, err := d.texts.Write(ctx, &proto.WriteTextRequest{Area: area, Name: user, Text: text})
  if err != nil {
    meta.Log().Fatal(fmt.Sprintf("text.Write failed: %v", err))
  }
}

func (d *Disk) ReadSquareCoords(area, user string) (*messages.Coords, error) {
  ctx, cancel := context.WithTimeout(context.Background(), d.timeout)
  defer cancel()

  resp, err := d.square.Read(ctx, &proto.ReadRequest{Area: area, Name: user})
  if err != nil {
    meta.Log().Warn(fmt.Sprintf("square.Read failed: %v", err))
    return nil, err
  }

  return &messages.Coords{XPos: resp.XPos, YPos: resp.YPos}, nil
}

func (d *Disk) ReadText(area, user string) (string, error) {
  ctx, cancel := context.WithTimeout(context.Background(), d.timeout)
  defer cancel()

  resp, err := d.texts.Read(ctx, &proto.ReadRequest{Area: area, Name: user})
  if err != nil {
    meta.Log().Warn(fmt.Sprintf("text.Read failed: %v", err))
    return "", err
  }

  return resp.Value, nil
}
