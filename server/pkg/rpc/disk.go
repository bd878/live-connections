package rpc

import (
  "time"
  "fmt"
  "context"

  "google.golang.org/grpc"

  "github.com/bd878/live-connections/disk/pkg/proto"

  "github.com/bd878/live-connections/meta"
  "github.com/bd878/live-connections/server/pkg/messages"
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

func (d *Disk) CreateNewUser(ctx context.Context, area string) (string, error) {
  ctx, cancel := context.WithTimeout(ctx, d.timeout)
  defer cancel()

  resp, err := d.user.Add(ctx, &proto.AddUserRequest{Area: area})
  if err != nil {
    meta.Log().Fatal(fmt.Sprintf("user.Add failed: %v", err))
    return "", err
  }

  return resp.Name, nil
}

func (d *Disk) CreateNewArea(ctx context.Context) (string, error) {
  ctx, cancel := context.WithTimeout(ctx, d.timeout)
  defer cancel()

  resp, err := d.area.Create(ctx, &proto.CreateAreaRequest{})
  if err != nil {
    meta.Log().Fatal(fmt.Sprintf("area.Create failed: %v", err))
    return "", err
  }

  return resp.Name, nil
}

func (d *Disk) ListUsers(ctx context.Context, area string) ([]string, error) {
  ctx, cancel := context.WithTimeout(ctx, d.timeout)
  defer cancel()

  resp, err := d.area.ListUsers(ctx, &proto.ListAreaUsersRequest{Name: area})
  if err != nil {
    meta.Log().Fatal(fmt.Sprintf("area.ListUsers failed: %v", err))
    return nil, err
  }

  return resp.GetUsers(), nil
}

func (d *Disk) HasUser(ctx context.Context, area, user string) (bool, error) {
  ctx, cancel := context.WithTimeout(ctx, d.timeout)
  defer cancel()

  resp, err := d.area.HasUser(ctx, &proto.HasUserRequest{Area: area, User: user})
  if err != nil {
    meta.Log().Fatal(fmt.Sprintf("area.HasUser failed: %v", err))
    return false, err
  }

  return resp.Result, nil
}

func (d *Disk) WriteSquareCoords(ctx context.Context, area, user string, XPos, YPos float32) error {
  ctx, cancel := context.WithTimeout(ctx, d.timeout)
  defer cancel()

  coords := &proto.Coords{XPos: XPos, YPos: YPos}
  _, err := d.square.Write(ctx, &proto.WriteSquareRequest{Area: area, Name: user, Coords: coords})
  if err != nil {
    meta.Log().Fatal(fmt.Sprintf("square.Write failed: %v", err))
    return err
  }

  return nil
}

func (d *Disk) WriteText(ctx context.Context, area, user string, recordId int32, value string) error {
  ctx, cancel := context.WithTimeout(ctx, d.timeout)
  defer cancel()

  text := &proto.Text{Value: value}
  _, err := d.texts.Write(ctx, &proto.WriteTextRequest{Area: area, Name: user, Text: text, RecordId: recordId})
  if err != nil {
    meta.Log().Fatal(fmt.Sprintf("text.Write failed: %v", err))
    return err
  }

  return nil
}

func (d *Disk) ReadSquareCoords(ctx context.Context, area, user string) (*messages.Coords, error) {
  ctx, cancel := context.WithTimeout(ctx, d.timeout)
  defer cancel()

  resp, err := d.square.Read(ctx, &proto.ReadRequest{Area: area, Name: user})
  if err != nil {
    meta.Log().Warn(fmt.Sprintf("square.Read failed: %v", err))
    return nil, err
  }

  return &messages.Coords{XPos: resp.XPos, YPos: resp.YPos}, nil
}

func (d *Disk) ReadText(ctx context.Context, area, user string, recordId int32) (string, error) {
  ctx, cancel := context.WithTimeout(ctx, d.timeout)
  defer cancel()

  resp, err := d.texts.Read(ctx, &proto.ReadRequest{Area: area, Name: user, RecordId: recordId})
  if err != nil {
    meta.Log().Warn(fmt.Sprintf("texts.ReadText failed: %v", err))
    return "", err
  }

  return resp.GetValue(), nil
}

func (d *Disk) AddTitle(ctx context.Context, area, user string) (*messages.Record, error) {
  ctx, cancel := context.WithTimeout(ctx, d.timeout)
  defer cancel()

  resp, err := d.texts.AddTitle(ctx, &proto.AddTitleRequest{Area: area, Name: user})
  if err != nil {
    meta.Log().Warn(fmt.Sprintf("texts.AddTitle failed: %v", err))
    return nil, err
  }

  meta.Log().Debug(
    fmt.Sprintf("Add title %d: %d, %d\n", resp.Id, resp.UpdatedAt, resp.CreatedAt),
  )

  result := &messages.Record{
    Value: resp.Value,
    ID: resp.Id,
    UpdatedAt: resp.UpdatedAt,
    CreatedAt: resp.CreatedAt,
  }

  return result, nil
}

func (d *Disk) ListTitles(ctx context.Context, area, user string) ([](*messages.Record), error) {
  ctx, cancel := context.WithTimeout(ctx, d.timeout)
  defer cancel()

  resp, err := d.texts.ListTitles(ctx, &proto.ListTitlesRequest{Area: area, Name: user})
  if err != nil {
    meta.Log().Warn(fmt.Sprintf("texts.ListTitles failed: %v", err))
    return nil, err
  }

  protoRecords := resp.GetRecords()
  result := make([](*messages.Record), len(protoRecords))

  for i, r := range protoRecords {
    result[i] = &messages.Record{
      Value: r.Value,
      ID: r.Id,
      CreatedAt: r.CreatedAt,
      UpdatedAt: r.UpdatedAt,
    }
  }

  return result, nil
}

func (d *Disk) ReadSelectedTitle(ctx context.Context, area, user string) (*messages.Record, error) {
  ctx, cancel := context.WithTimeout(ctx, d.timeout)
  defer cancel()

  resp, err := d.texts.ReadSelectedTitle(ctx, &proto.ReadSelectedRequest{Area: area, Name: user})
  if err != nil {
    meta.Log().Warn(fmt.Sprintf("texts.ReadSelectedTitle failed: %v", err))
    return nil, err
  }

  return &messages.Record{
    Value: resp.Value,
    CreatedAt: resp.CreatedAt,
    UpdatedAt: resp.UpdatedAt,
  }, nil
}

func (d *Disk) SelectTitle(ctx context.Context, area, user string, recordId int32) error {
  ctx, cancel := context.WithTimeout(ctx, d.timeout)
  defer cancel()

  _, err := d.texts.SelectTitle(ctx, &proto.SelectTitleRequest{Area: area, Name: user, RecordId: recordId})
  if err != nil {
    meta.Log().Warn(fmt.Sprintf("texts.SelectTitle failed: %v", err))
    return err
  }

  return nil
}
