package services

import (
  "errors"
  "log"
  "fmt"
  "time"
  "strings"
  "context"
  "regexp"
  "bufio"
  "os"
  "path/filepath"
  "encoding/binary"

  "google.golang.org/protobuf/proto"

  "github.com/bd878/live-connections/disk/pkg/utils"
  pb "github.com/bd878/live-connections/disk/pkg/proto"
)

type TextsServer struct {
  pb.UnimplementedTextsManagerServer
  Dir string
}

func buildFilename(id int32) string {
  var b strings.Builder
  fmt.Fprintf(&b, "%d.text.state", id)
  return b.String()
}

var textRe = regexp.MustCompile(`^\d*\.text\.state`)

const selectedTextName = "selected.state"

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

  fp := filepath.Join(s.Dir, request.Area, request.Name, buildFilename(request.RecordId))

  log.Println("write coords in file =", fp)

  f, err := os.OpenFile(
    fp,
    os.O_RDWR|os.O_CREATE,
    0644,
  )
  if err != nil {
    log.Println("failed to open file =", fp)
    return nil, err
  }

  info, err := f.Stat()
  if err != nil {
    return nil, err
  }

  size := info.Size()
  p := make([]byte, size)

  log.Println("file size:", size)

  bytesRead, err := f.Read(p)
  if err != nil {
    return nil, err
  }

  log.Println("bytes read:", bytesRead)

  if int64(bytesRead) != size {
    return nil, errors.New("read less bytes than in file")
  }

  title := &pb.TitleRecord{}
  text := &pb.Text{}
  textRecord := &pb.TextRecord{Title: title, Text: text}
  if len(p) > 0 {
    if err = proto.Unmarshal(p, textRecord); err != nil {
      return nil, err
    }
  }

  textRecord.Text.Value = request.Text.Value

  log.Println("text value:", request.Text.Value)

  p, err = proto.Marshal(textRecord)
  if err != nil {
    log.Println("failed to marshal text record")
    return nil, err
  }

  buf := bufio.NewWriter(f)
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
  // TODO: look up for same logic (listTitles), unify
  if !utils.IsNameSafe(request.Area) {
    return nil, errors.New("area name not safe")
  }

  if !utils.IsNameSafe(request.Name) {
    return nil, errors.New("user name not safe")
  }

  fp := filepath.Join(s.Dir, request.Area, request.Name, buildFilename(request.RecordId))
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

  title := &pb.TitleRecord{}
  text := &pb.Text{}
  textRecord := &pb.TextRecord{Title: title, Text: text}
  if err = proto.Unmarshal(p, textRecord); err != nil {
    return nil, err
  }

  return textRecord.Text, nil
}

func (s *TextsServer) AddTitle(ctx context.Context, request *pb.AddTitleRequest) (*pb.TitleRecord, error) {
  log.Println("add title")

  if !utils.IsNameSafe(request.Area) {
    return nil, errors.New("area name not safe")
  }

  if !utils.IsNameSafe(request.Name) {
    return nil, errors.New("user name not safe")
  }

  createdAt := int32(time.Now().Unix())
  updatedAt := createdAt

  fmt.Println("name, createdAt, updatedAt", request.Name, createdAt, updatedAt)

  fp := filepath.Join(s.Dir, request.Area, request.Name, buildFilename(createdAt))
  f, err := os.OpenFile(
    fp,
    os.O_WRONLY|os.O_CREATE|os.O_TRUNC,
    0644,
  )
  if err != nil {
    return nil, err
  }

  title := &pb.TitleRecord{
    Value: "",
    Id: createdAt,
    CreatedAt: createdAt,
    UpdatedAt: updatedAt,
  }
  text := &pb.Text{Value: ""}
  result := &pb.TextRecord{
    Title: title,
    Text: text,
  }

  p, err := proto.Marshal(result)
  if err != nil {
    log.Println("failed to marshal text record")
    return nil, err
  }

  buf := bufio.NewWriter(f)
  if _, err = buf.Write(p); err != nil {
    log.Printf("error write to buffer = %v\n", err)
    return nil, err
  }

  if err = buf.Flush(); err != nil {
    log.Printf("error flush to file = %v\n", err)
    return nil, err
  }

  return title, nil
}

func (s *TextsServer) ListTitles(ctx context.Context, request *pb.ListTitlesRequest) (*pb.ListTitlesResponse, error) {
  log.Println("list titles")

  if !utils.IsNameSafe(request.Area) {
    return nil, errors.New("area name not safe")
  }

  if !utils.IsNameSafe(request.Name) {
    return nil, errors.New("user name not safe")
  }

  dp := filepath.Join(s.Dir, request.Area, request.Name)
  files, err := os.ReadDir(dp)
  if err != nil {
    return nil, err
  }

  records := make([]*pb.TitleRecord, 0, 100)
  for _, file := range files {
    if textRe.MatchString(file.Name()) {
      // TODO: unify with Read, same logic
      fp := filepath.Join(s.Dir, request.Area, request.Name, file.Name())

      f, err := os.OpenFile(
        fp,
        os.O_RDONLY,
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

      title := &pb.TitleRecord{}
      text := &pb.Text{}
      textRecord := &pb.TextRecord{Title: title, Text: text}
      if err = proto.Unmarshal(p, textRecord); err != nil {
        return nil, err
      }

      records = append(records, textRecord.Title)
    }
  }

  result := &pb.ListTitlesResponse{
    Records: records,
  }

  return result, nil
}

func (s *TextsServer) SelectTitle(ctx context.Context, request *pb.SelectTitleRequest) (*pb.EmptyResponse, error) {
  log.Println("select title")

  // TODO: check if record file exists
  if !utils.IsNameSafe(request.Area) {
    return nil, errors.New("area name not safe")
  }

  if !utils.IsNameSafe(request.Name) {
    return nil, errors.New("user name not safe")
  }

  fp := filepath.Join(s.Dir, request.Area, request.Name, selectedTextName)
  f, err := os.OpenFile(
    fp,
    os.O_WRONLY|os.O_CREATE|os.O_TRUNC,
    0644,
  )
  if err != nil {
    log.Println("failed to open file =", fp)
    return nil, err
  }

  err = binary.Write(f, binary.LittleEndian, request.RecordId)
  if err != nil {
    return nil, err
  }

  return &pb.EmptyResponse{}, nil
}

func (s *TextsServer) ReadSelectedTitle(ctx context.Context, request *pb.ReadSelectedRequest) (*pb.TitleRecord, error) {
  log.Println("read selected title")

  // TODO: check if record file exists
  if !utils.IsNameSafe(request.Area) {
    return nil, errors.New("area name not safe")
  }

  if !utils.IsNameSafe(request.Name) {
    return nil, errors.New("user name not safe")
  }

  fp := filepath.Join(s.Dir, request.Area, request.Name, selectedTextName)
  f, err := os.OpenFile(
    fp,
    os.O_RDONLY,
    0644,
  )
  if err != nil {
    log.Println("failed to open file =", fp)
    return nil, err
  }

  log.Println("selected filepath:", fp)

  var recordId int32
  err = binary.Read(f, binary.LittleEndian, &recordId)
  if err != nil {
    return nil, err
  }

  log.Println("selected recordId", recordId)

  recordFp := filepath.Join(s.Dir, request.Area, request.Name, buildFilename(recordId))

  log.Println("record fp:", recordFp)

  recordF, err := os.OpenFile(
    recordFp,
    os.O_RDONLY,
    0644,
  )
  if err != nil {
    return nil, err
  }

  info, err := recordF.Stat()
  if err != nil {
    return nil, err
  }

  size := info.Size()
  p := make([]byte, size)

  bytesRead, err := recordF.Read(p)
  if err != nil {
    return nil, err
  }

  log.Println("size / bytes read:", size, bytesRead)

  if int64(bytesRead) != size {
    return nil, errors.New("read less bytes than in file")
  }

  title := &pb.TitleRecord{}
  text := &pb.Text{}
  textRecord := &pb.TextRecord{Title: title, Text: text}
  if err = proto.Unmarshal(p, textRecord); err != nil {
    return nil, err
  }

  log.Println("text record title:", textRecord.Title)

  return textRecord.Title, nil
}