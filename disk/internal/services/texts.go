package services

import (
  "errors"
  "log"
  "context"
  "bufio"
  "os"
  "path/filepath"

  "google.golang.org/protobuf/proto"

  "github.com/teralion/live-connections/disk/pkg/utils"
  pb "github.com/teralion/live-connections/disk/pkg/proto"
)

const textsFileName = "text.state"

type TextsServer struct {
  pb.UnimplementedTextsManagerServer
  Dir string
}

func NewTextsManagerServer(baseDir string) *TextsServer {
  return &TextsServer{Dir: baseDir}
}

func (s *TextsServer) Write(ctx context.Context, request *pb.WriteTextRequest) (*pb.EmptyResponse, error) {
  log.Println("write text =", request.Text.Value)

  if !utils.IsNameSafe(request.Area) {
    return nil, errors.New("area name not safe")
  }

  if !utils.IsNameSafe(request.Name) {
    return nil, errors.New("user name not safe")
  }

  fp := filepath.Join(s.Dir, request.Area, request.Name, textsFileName)

  log.Println("write coords in file =", fp)

  p, err := proto.Marshal(request.Text)
  if err != nil {
    log.Println("failed to marshal text")
    return nil, err
  }

  storeFile, err := os.OpenFile(
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

  return &pb.EmptyResponse{}, nil
}

func (s *TextsServer) Read(ctx context.Context, request *pb.ReadRequest) (*pb.Text, error) {
  if !utils.IsNameSafe(request.Area) {
    return nil, errors.New("area name not safe")
  }

  if !utils.IsNameSafe(request.Name) {
    return nil, errors.New("user name not safe")
  }

  fp := filepath.Join(s.Dir, request.Area, request.Name, textsFileName)
  f, err := os.OpenFile(
    fp,
    os.O_RDONLY|os.O_CREATE,
    0644,
  )
  if err != nil {
    return nil, err
  }

  info, err := f.Stat()
  if err != nil {
    return nil, err
  }

  size := info.Size()
  p := make([]byte, size)

  bytesRead, err := f.Read(p)
  if err != nil {
    return nil, err
  }

  if int64(bytesRead) != size {
    return nil, errors.New("read less bytes than in file")
  }

  text := &pb.Text{}
  if err = proto.Unmarshal(p, text); err != nil {
    return nil, err
  }

  return text, nil
}