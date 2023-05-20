package services

import (
  "fmt"
  "time"
  "strings"
  "context"
  "regexp"
  "bufio"
  "os"
  "encoding/binary"

  "google.golang.org/protobuf/proto"

  "github.com/bd878/live-connections/meta"
  "github.com/bd878/live-connections/disk/pkg/fd"

  pb "github.com/bd878/live-connections/disk/pkg/proto"
)

func buildFilename(id int32) string {
  var b strings.Builder
  fmt.Fprintf(&b, "%d.text.state", id)
  return b.String()
}

type TextsServer struct {
  pb.UnimplementedTextsManagerServer
  Dir string
}

var textRe = regexp.MustCompile(`^\d*\.text\.state`)

const selectedTextName = "selected.state"

func NewTextsManagerServer(baseDir string) *TextsServer {
  return &TextsServer{Dir: baseDir}
}

func (s *TextsServer) Write(ctx context.Context, request *pb.WriteTextRequest) (*pb.EmptyResponse, error) {
  f := fd.NewFile(os.O_RDWR|os.O_CREATE)

  err := f.Open(s.Dir, request.Area, request.Name, buildFilename(request.RecordId))
  if err != nil {
    return nil, err
  }

  err = f.Load()
  if err != nil {
    return nil, err
  }

  textRecord := &pb.TextRecord{}
  if err = proto.Unmarshal(f.Content(), textRecord); err != nil {
    return nil, err
  }

  textRecord.Text.Value = request.Text.Value

  p, err := proto.Marshal(textRecord)
  if err != nil {
    return nil, err
  }

  buf := bufio.NewWriter(f.File())
  if _, err = buf.Write(p); err != nil {
    meta.Log().Debug(fmt.Sprintf("error write to buffer = %v\n", err))
    return nil, err
  }

  if err = buf.Flush(); err != nil {
    meta.Log().Debug(fmt.Sprintf("error flush to file = %v\n", err))
    return nil, err
  }

  return &pb.EmptyResponse{}, nil
}

func (s *TextsServer) Read(ctx context.Context, request *pb.ReadRequest) (*pb.Text, error) {
  f := fd.NewFile(os.O_RDONLY|os.O_CREATE)

  err := f.Open(s.Dir, request.Area, request.Name, buildFilename(request.RecordId))
  if err != nil {
    return nil, err
  }

  err = f.Load()
  if err != nil {
    return nil, err
  }

  textRecord := &pb.TextRecord{}
  if err = proto.Unmarshal(f.Content(), textRecord); err != nil {
    return nil, err
  }

  result := &pb.Text{
    Value: textRecord.Text.Value,
  }

  return result, nil
}

func (s *TextsServer) Add(ctx context.Context, request *pb.AddTextRecordRequest) (*pb.TextRecord, error) {
  createdAt := int32(time.Now().Unix())
  updatedAt := createdAt

  result := &pb.TextRecord{
    Text: &pb.Text{
      Value: "",
    },
    Title: "",
    Id: createdAt,
    CreatedAt: createdAt,
    UpdatedAt: updatedAt,
  }

  p, err := proto.Marshal(result)
  if err != nil {
    meta.Log().Debug("failed to marshal text record")
    return nil, err
  }

  f := fd.NewFile(os.O_WRONLY|os.O_CREATE|os.O_TRUNC)

  err = f.Open(s.Dir, request.Area, request.Name, buildFilename(result.Id))
  if err != nil {
    return nil, err
  }

  err = f.Load()
  if err != nil {
    return nil, err
  }

  buf := bufio.NewWriter(f.File())
  if _, err = buf.Write(p); err != nil {
    meta.Log().Debug(fmt.Sprintf("error write to buffer = %v\n", err))
    return nil, err
  }

  if err = buf.Flush(); err != nil {
    meta.Log().Debug(fmt.Sprintf("error flush to file = %v\n", err))
    return nil, err
  }

  return result, nil
}

func (s *TextsServer) List(ctx context.Context, request *pb.ListTextRecordsRequest) (*pb.ListTextRecordsResponse, error) {
  f := fd.NewDir()

  err := f.Open(s.Dir, request.Area, request.Name)
  if err != nil {
    return nil, err
  }

  records := make([]*pb.TextRecord, 0, 100)
  for _, file := range f.Content() {
    if textRe.MatchString(file.Name()) {
      f := fd.NewFile(os.O_RDONLY)

      err := f.Open(s.Dir, request.Area, request.Name, file.Name())
      if err != nil {
        return nil, err
      }

      err = f.Load()
      if err != nil {
        return nil, err
      }

      textRecord := &pb.TextRecord{}
      if err = proto.Unmarshal(f.Content(), textRecord); err != nil {
        return nil, err
      }

      records = append(records, textRecord)
    }
  }

  result := &pb.ListTextRecordsResponse{
    Records: records,
  }

  return result, nil
}

func (s *TextsServer) Select(ctx context.Context, request *pb.SelectTextRecordRequest) (*pb.EmptyResponse, error) {
  f := fd.NewFile(os.O_WRONLY|os.O_CREATE|os.O_TRUNC)

  err := f.Open(s.Dir, request.Area, request.Name, selectedTextName)
  if err != nil {
    return nil, err
  }

  err = binary.Write(f.File(), binary.LittleEndian, request.RecordId)
  if err != nil {
    return nil, err
  }

  return &pb.EmptyResponse{}, nil
}

func (s *TextsServer) GetSelected(ctx context.Context, request *pb.GetSelectedRecordRequest) (*pb.TextRecord, error) {
  f := fd.NewFile(os.O_RDONLY)

  err := f.Open(s.Dir, request.Area, request.Name, selectedTextName)
  if err != nil {
    return nil, err
  }

  var recordId int32
  err = binary.Read(f.File(), binary.LittleEndian, &recordId)
  if err != nil {
    return nil, err
  }

  recordF := fd.NewFile(os.O_RDONLY)
  err = recordF.Open(s.Dir, request.Area, request.Name, buildFilename(recordId))
  if err != nil {
    return nil, err
  }

  err = recordF.Load()
  if err != nil {
    return nil, err
  }

  textRecord := &pb.TextRecord{}
  if err = proto.Unmarshal(recordF.Content(), textRecord); err != nil {
    return nil, err
  }

  return textRecord, nil
}