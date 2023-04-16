package main

import (
  "fmt"
  "log"
  "time"
)

type Server interface {
  ListenAndServe()
  Shutdown()
}

type PipeBackend interface {
  Respond() string
  LoadZones(p string) string
}

type BackendServer struct{}

func New() Server {
  return &BackendServer{}
}

func (s *BackendServer) ListenAndServe() {
  fmt.Println("listening...")
}

func (s *BackendServer) Shutdown() {
  fmt.Println("shutdown")
}

func (s *BackendServer) Respond() string {
  fmt.Println("responding")
  return "tadam"
}

func (s *BackendServer) LoadZones(p string) string {
  return p
}

func main() {
  s := New()
  s.ListenAndServe()
  time.Sleep(1 * time.Second)
  p, ok := s.(PipeBackend)
  if !ok {
    log.Fatal("not a pipe backend")
  }
  p.Respond()
  s.Shutdown()
}
