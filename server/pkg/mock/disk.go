package mock

import (
  "errors"
  "context"
  "time"

  "github.com/bd878/live-connections/server/pkg/messages"
  "github.com/bd878/live-connections/disk/pkg/utils"
)

type Disk struct {
  users []string
  areas []string
  squareCoords map[string]*messages.Coords
  records map[string][]*messages.Record
  selectedRecord *messages.Record
}

func NewDisk() *Disk {
  return &Disk{}
}

func (d *Disk) CreateNewUser(ctx context.Context, area string) (string, error) {
  user := utils.RandomString(10)
  d.users = append(d.users, user)
  return user, nil
}

func (d *Disk) CreateNewArea(ctx context.Context) (string, error) {
  area := utils.RandomString(10)
  d.areas = append(d.areas, area)
  return area, nil
}

func (d *Disk) ListUsers(ctx context.Context, area string) ([]string, error) {
  return d.users, nil
}

func (d *Disk) HasUser(ctx context.Context, _ string, user string) (bool, error)   {
  var found bool
  for i := 0; i < len(d.users) && !found; i++ {
    if d.users[i] == user {
      found = true
    }
  }
  return found, nil
}


func (d *Disk) WriteSquareCoords(ctx context.Context, _ string, user string, XPos, YPos float32) error {
  d.squareCoords[user] = &messages.Coords{XPos: XPos, YPos: YPos}
  return nil
}

func (d *Disk) ReadSquareCoords(ctx context.Context, _ string, user string) (*messages.Coords, error) {
  return d.squareCoords[user], nil
}

func (d *Disk) findRecordById(ctx context.Context, user string, recordId int32) (*messages.Record, error) {
  recs := d.records[user]

  var result *messages.Record

  var found bool
  for i := 0; i < len(recs) && !found; i++ {
    rec := recs[i]
    if rec.ID == recordId {
      found = true
      result = rec
    }
  }

  if found {
    return result, nil
  } else {
    return nil, errors.New("not found")
  }
}

func (d *Disk) WriteText(ctx context.Context, _ string, user string, recordId int32, value string) error {
  rec, err := d.findRecordById(ctx, user, recordId)
  if err != nil {
    return err
  }

  rec.Value = value
  return nil
}

func (d *Disk) ReadText(ctx context.Context, _ string, user string, recordId int32) (string, error) {
  rec, err := d.findRecordById(ctx, user, recordId)
  if err != nil {
    return "", err
  }
  return rec.Value, nil
}


func (d *Disk) AddTextRecord(ctx context.Context, _ string, user string) (*messages.Record, error) {
  updatedAt := int32(time.Now().Unix())
  createdAt := updatedAt
  id := createdAt

  rec := &messages.Record{
    Value: "",
    ID: id,
    UpdatedAt: updatedAt,
    CreatedAt: createdAt,
  }

  d.records[user] = append(d.records[user], rec)
  return rec, nil
}

func (d *Disk) ListTextRecords(ctx context.Context, _ string, user string) ([](*messages.Record), error) {
  return d.records[user], nil
}

func (d *Disk) GetSelectedRecord(ctx context.Context, _ string, user string) (*messages.Record, error) {
  if d.selectedRecord != nil {
    return d.selectedRecord, nil
  }
  return nil, errors.New("not selected")
}

func (d *Disk) SelectTextRecord(ctx context.Context, _ string, user string, recordId int32) error {
  rec, err := d.findRecordById(ctx, user, recordId)
  if err != nil {
    return err
  }
  d.selectedRecord = rec
  return nil
}
