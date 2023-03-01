package services

import (
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

const squareFileName = "square.state"

type SquareServer struct {
  pb.UnimplementedSquareManagerServer
  Dir string
}

func NewSquareManagerServer(baseDir string) *SquareServer {
  return &SquareServer{Dir: baseDir}
}

func (s *SquareServer) Write(ctx context.Context, request *pb.WriteSquareRequest) (*pb.WriteSquareResponse, error) {
  var (
    p []byte
    err error
    storeFile *os.File
  )

  log.Println("write mouse coords =", request.Coords.XPos, request.Coords.YPos)

  if !utils.IsNameSafe(request.Area) {
    return nil, errors.New("area name not safe")
  }

  if !utils.IsNameSafe(request.Name) {
    return nil, errors.New("user name not safe")
  }

  fp := filepath.Join(s.Dir, request.Area, request.Name, squareFileName)

  log.Println("write coords in file =", fp)

  if p, err = proto.Marshal(request.Coords); err != nil {
    log.Println("failed to marshal coords")
    return nil, err
  }

  storeFile, err = os.OpenFile(
    fp,
    os.O_WRONLY|os.O_CREATE|os.O_TRUNC,
    0644,
  )
  if err != nil {
    log.Println("failed to open file =", fp)
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

  return &pb.WriteSquareResponse{}, nil
}

func (s *SquareServer) Read(ctx context.Context, request *pb.ReadSquareRequest) (*pb.Coords, error) {
  var (
    err error
    f *os.File
    size int64
    info os.FileInfo
    bytesRead int
  )

  if !utils.IsNameSafe(request.Area) {
    return nil, errors.New("area name not safe")
  }

  if !utils.IsNameSafe(request.Name) {
    return nil, errors.New("user name not safe")
  }

  fp := filepath.Join(s.Dir, request.Area, request.Name, squareFileName)
  f, err = os.OpenFile(
    fp,
    os.O_RDONLY|os.O_CREATE,
    0644,
  )
  if err != nil {
    return nil, err
  }

  if info, err = f.Stat(); err != nil {
    return nil, err
  }

  size = info.Size()
  p := make([]byte, size)

  bytesRead, err = f.Read(p)
  if err != nil {
    return nil, err
  }

  if int64(bytesRead) != size {
    return nil, errors.New("read less bytes than in file")
  }

  coords := &pb.Coords{}
  if err = proto.Unmarshal(p, coords); err != nil {
    return nil, err
  }

  return coords, nil
}

