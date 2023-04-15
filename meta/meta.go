package meta

import (
  "log"
  "sync"
  "os"
)

type LogWriter struct {
  prefix string
  degree string
}

var (
  once sync.Once
  instance LogWriter
)

func NewLogWriter(prefix string, degree string) *LogWriter {
  once.Do(func() {
    if degree == "" {
      degree = "warn"
    }

    instance = LogWriter{
      prefix: prefix,
      degree: degree,
    }
  })

  return &instance
}

func Log() *LogWriter {
  return NewLogWriter("", "debug")
}

func (l *LogWriter) Info(text ...any) {
  if l.degree != "silent" {
    log.Println(text...)
  }
}

func (l *LogWriter) Debug(text ...any) {
  if l.degree != "silent" && (l.degree == "debug") {
    log.Println(text...)
  }
}

func (l *LogWriter) Warn(text ...any) {
  if l.degree != "silent" && (l.degree == "debug" || l.degree == "warn") {
    log.Println(text...)
  }
}

func (l *LogWriter) Error(text ...any) {
  if l.degree != "silent" && (l.degree == "debug" || l.degree == "warn" || l.degree == "fatal") {
    log.Println(text...)
  }
}

func (l *LogWriter) Fatal(text ...any) {
  l.Error(text...)
  os.Exit(1)
}
