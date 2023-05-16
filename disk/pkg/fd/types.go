package fd

import "os"

type Node interface {
  Open(paths ...string) error
}

type FileNode interface {
  SetFlags(flags int)
  Content() []byte
  Load() error
}

type DirNode interface {
  Content() []os.DirEntry
}