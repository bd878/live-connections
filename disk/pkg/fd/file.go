package fd

import (
  "os"
  "fmt"
  "errors"
  "path/filepath"

  "github.com/bd878/live-connections/disk/pkg/utils"
)

func isSafePaths(paths ...string) error {
  for _, p := range paths {
    if !utils.IsNameSafe(p) {
      return errors.New(fmt.Sprintf("%s not safe", p))
    }
  }

  return nil
}

type File struct {
  file *os.File
  data []byte
  flags int
}

func NewFile(flags int) *File {
  return &File{
    flags: flags,
  }
}

func (tf *File) Open(paths ...string) error {
  err := isSafePaths(paths...)
  if err != nil {
    return err
  }

  f, err := os.OpenFile(
    filepath.Join(paths...),
    tf.flags,
    0644,
  )
  if err != nil {
    return err
  }

  tf.file = f
  return nil
}

func (tf *File) File() *os.File {
  return tf.file
}

func (tf *File) Load() error {
  info, err := tf.File().Stat()
  if err != nil {
    return err
  }

  size := info.Size()
  tf.data = make([]byte, size)

  bytesRead, err := tf.File().Read(tf.data)
  if err != nil {
    return err
  }

  if int64(bytesRead) != size {
    return errors.New("read less bytes than in file")
  }

  return nil
}

func (tf *File) SetFlags(flags int) {
  tf.flags = flags
}

func (tf *File) Content() []byte {
  return tf.data
}
