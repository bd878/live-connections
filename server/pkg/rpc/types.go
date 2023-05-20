package rpc

import (
  "context"

  "github.com/bd878/live-connections/server/pkg/messages"
)

type IdentityDisk interface {
  CreateNewUser(ctx context.Context, area string) (string, error)
  CreateNewArea(ctx context.Context) (string, error)
  ListUsers(ctx context.Context, area string) ([]string, error)
  HasUser(ctx context.Context, area, user string) (bool, error)  
}

type SquareDisk interface {
  WriteSquareCoords(ctx context.Context, area, user string, XPos, YPos float32) error
  ReadSquareCoords(ctx context.Context, area, user string) (*messages.Coords, error)
}

type TextDisk interface {
  WriteText(ctx context.Context, area, user string, recordId int32, value string) error
  ReadText(ctx context.Context, area, user string, recordId int32) (string, error)
}

type CatalogDisk interface {
  AddTextRecord(ctx context.Context, area, user string) (*messages.TextRecord, error)
  ListTextRecords(ctx context.Context, area, user string) ([](*messages.TextRecord), error)
  GetSelectedRecord(ctx context.Context, area, user string) (*messages.TextRecord, error)
  SelectTextRecord(ctx context.Context, area, user string, recordId int32) error
}
