package services

import (
  "os"
  "testing"
  "context"
  "fmt"
  "reflect"
  "path/filepath"

  pb "github.com/teralion/live-connections/disk/pkg/proto"
)

const userNameLength = 10

type Config struct {
  Area string
}

func TestUser(t *testing.T) {
  for scenario, fn := range map[string]func(
    t *testing.T,
    user *UserManagerServer,
    config *Config,
  ) {
    "create user": testUserCreate,
  } {
    t.Run(scenario, func(t *testing.T) {
      user, config := setupTestUser(t)
      fn(t, user, config)
    })
  }
}

func setupTestUser(t *testing.T) (*UserManagerServer, *Config) {
  t.Helper()

  dir := os.TempDir()
  u := UserManagerServer{Dir: dir, NameLen: userNameLength}

  a := AreaManagerServer{Dir: dir, NameLen: 10}
  respArea, _ := a.Create(context.Background(), &pb.CreateAreaRequest{})
  config := Config{Area: respArea.Name}

  return &u, &config
}

func testUserCreate(t *testing.T, user *UserManagerServer, config *Config) {
  resp, err := user.Add(context.Background(), &pb.AddUserRequest{Area: config.Area})
  if err != nil {
    t.Fatal(err)
  }

  if _, found := reflect.TypeOf(*resp).FieldByName("Name"); !found {
    t.Fatal("AddUser has no 'Name' field")
  }

  if len(resp.Name) != userNameLength {
    fmt.Errorf("name length != required length : %d != %d", len(resp.Name), userNameLength)
  }

  if _, err := os.Open(filepath.Join(user.Dir, config.Area, resp.Name)); err != nil {
    t.Fatal(err)
  }
}