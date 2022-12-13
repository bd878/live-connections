package services

import (
  "testing"
  "context"
  "reflect"
  "fmt"
  "os"

  pb "github.com/teralion/live-connections/disk/pkg/proto"
)

const areaNameLen = 10

func TestArea(t *testing.T) {
  for scenario, fn := range map[string]func(
    t *testing.T,
    area *AreaManagerServer,
  ) {
    "test create": testCreate,
  } {
    t.Run(scenario, func(t *testing.T) {
      area := setupTest(t)
      fn(t, area)
    })
  }
}

func setupTest(t *testing.T) *AreaManagerServer {
  t.Helper()

  dir := os.TempDir()
  a := AreaManagerServer{Dir: dir, NameLen: areaNameLen}

  return &a
}

func testCreate(t *testing.T, area *AreaManagerServer) {
  resp, err := area.Create(context.Background(), &pb.CreateAreaRequest{})
  if err != nil {
    t.Fatal(err)
  }

  if _, found := reflect.TypeOf(*resp).FieldByName("Name"); !found {
    t.Fatal("has no 'Name' field")
  }

  if len(resp.Name) != areaNameLen {
    fmt.Errorf("name length != required length : %d != %d", len(resp.Name), areaNameLen)
  }
}