package meta

import (
  "log"
  "sync"
  "os"
)

type LogWriter struct {
  degree string
}

var (
  once sync.Once
  instance LogWriter
)

func NewLogWriter(degree string) *LogWriter {
  once.Do(func() {
    if degree == "" {
      degree = "debug"
    }

    instance = LogWriter{degree}
  })

  return &instance
}

func Log() *LogWriter {
  return NewLogWriter("debug")
}

func (l *LogWriter) Info(text ...any) {
  log.Println(text...)
}

func (l *LogWriter) Debug(text ...any) {
  log.Println(text...)
}

func (l *LogWriter) Warn(text ...any) {
  log.Println(text...)
}

func (l *LogWriter) Error(text ...any) {
  log.Println(text...)
}

func (l *LogWriter) Fatal(text ...any) {
  l.Error(text...)
  os.Exit(1)
}
