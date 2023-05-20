package services

import (
  "os"
  "path/filepath"
  "testing"
  "context"

  "github.com/bd878/live-connections/disk/pkg/utils"
  pb "github.com/bd878/live-connections/disk/pkg/proto"
)

func TestTexts(t *testing.T) {
  area := utils.RandomString(10)
  user := utils.RandomString(10)

  dir := filepath.Join("./testdata", area, user)
  err := os.MkdirAll(dir, 0750)
  if err != nil && !os.IsExist(err) {
    t.Fatal(err)
  }

  server := NewTextsManagerServer("./testdata")

  addRequest := &pb.AddTextRecordRequest{Area: area, Name: user}
  record, err := server.Add(context.TODO(), addRequest)
  if err != nil {
    t.Fatal(err)
  }

  t.Log(record.Id)
}
