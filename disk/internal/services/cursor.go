package services

import (
  "fmt"
  "errors"
  "log"
  "context"
  "bufio"
  "os"
  "path/filepath"

  "google.golang.org/protobuf/proto"

  utils "github.com/teralion/live-connections/disk/internal/utils"
  pb "github.com/teralion/live-connections/disk/pkg/proto"
)

const cursorFileName = "cursor.state"

type CursorServer struct {
  pb.UnimplementedCursorManagerServer
  Dir string
}

func NewCursorManagerServer(baseDir string) *CursorServer {
  return &CursorServer{Dir: baseDir}
}

func (s *CursorServer) Write(ctx context.Context, request *pb.WriteCursorRequest) (*pb.WriteCursorResponse, error) {
  var (
    p []byte
    err error
    storeFile *os.File
  )

  if !utils.IsNameSafe(request.Area) {
    return nil, errors.New("area name not safe")
  }

  if !utils.IsNameSafe(request.Name) {
    return nil, errors.New("user name not safe")
  }

  fp := filepath.Join(s.Dir, request.Area, request.Name, cursorFileName)

  if p, err = proto.Marshal(request.Coords); err != nil {
    return nil, err
  }

  storeFile, err = os.OpenFile(
    fp,
    os.O_WRONLY|os.O_CREATE|os.O_TRUNC,
    0644,
  )
  if err != nil {
    return nil, err
  }

  buf := bufio.NewWriter(storeFile)
  if _, err = buf.Write(p); err != nil {
    log.Printf("error write to buffer = %v\n", err)
    return nil, err
  }

  if err = buf.Flush(); err != nil {
    log.Printf("error flush to file = %v\n", err)
    return nil, err
  }

  return &pb.WriteCursorResponse{}, nil
}

func (s *CursorServer) Read(ctx context.Context, request *pb.ReadCursorRequest) (*pb.ReadCursorResponse, error) {
  fmt.Println("not implemented")
  return &pb.ReadCursorResponse{XPos: 0, YPos: 0}, nil
}

