package protocol

import (
  "testing"
  "reflect"
  "time"
  "context"

  "github.com/bd878/live-connections/server/pkg/mock"
)

func TestClientName(t *testing.T) {
  client := NewClient(mock.NewConn())

  name := "test"

  client.SetName(name)

  if client.Name() != name {
    t.Fatal("client returned wrong name", name)
  }
}

func TestClientParent(t *testing.T) {
  area1 := mock.NewArea()
  client := NewClient(mock.NewConn())

  client.SetParent(area1)

  area2 := client.Parent()

  pValue1 := reflect.ValueOf(area1)
  pValue2 := reflect.ValueOf(area2)

  switch pValue2.Kind() {
  case reflect.Pointer:
    t.Log("area2 is a pointer")
  default:
    t.Fatal("area2 is NOT a pointer")
  }

  if !pValue1.Equal(pValue2) {
    t.Fatal("area1 is NOT equal area2")
  } else {
    t.Log("area1 equal area2")
  }
}

func TestClientSend(t *testing.T) {
  conn := mock.NewConn()
  client := NewClient(conn)

  go client.Run(context.Background())
  time.Sleep(100*time.Millisecond)

  var d = []byte{1, 2, 3}

  client.Send() <- d
  result := <-conn.Messages()

  d = append(d, []byte{'\n'}...)

  if !reflect.DeepEqual(d, result) {
    t.Fatal("conn received different value than was sent")
  }
}

func TestClientDisk(t *testing.T) {
  disk1 := mock.NewDisk()
  client := NewClient(mock.NewConn())

  client.SetDisk(disk1)

  disk2 := client.Disk()

  dValue1 := reflect.ValueOf(disk1)
  dValue2 := reflect.ValueOf(disk2)

  if !dValue1.Equal(dValue2) {
    t.Fatal("disk1 is NOT equal disk2")
  } else {
    t.Log("disk1 equal disk2")
  }
}

func TestClientRestore(t *testing.T) {
  disk := mock.NewDisk()
  area := mock.NewArea()
  conn := mock.NewConn()

  client := NewClient(conn)
  client.SetParent(area)
  client.SetDisk(disk)

  areaName := "area"
  clientName := "client"
  inputText := "testtest"
  squareXPos = 30
  squareYPos = 20

  area.SetName(areaName)
  client.SetSquareX(squareXPos)
  client.SetSquareY(squareYPos)
  client.SetName(clientName)
  client.SelectRecord(client.NewRecord())
  client.SetInputText(inputText)

  client.Save()
}