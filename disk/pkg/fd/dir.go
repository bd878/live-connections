package fd

import (
  "os"
  "fmt"
  "errors"
  "path/filepath"

  "github.com/bd878/live-connections/disk/pkg/utils"
)

type Dir struct {
  entries []os.DirEntry
}

func NewDir() *Dir {
  return &Dir{}
}

func (d *Dir) Open(paths ...string) error {
  for _, p := range paths {
    if !utils.IsNameSafe(p) {
      return errors.New(fmt.Sprintf("%s not safe", p))
    }
  }

  ds, err := os.ReadDir(filepath.Join(paths...))
  if err != nil {
    return err
  }

  d.entries = ds
  return nil
}

func (d *Dir) Content() []os.DirEntry {
  return d.entries
}
