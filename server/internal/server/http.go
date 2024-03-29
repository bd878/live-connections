package server

import (
  "net/http"
  "fmt"
  "time"
  "os"
  "context"
  "os/signal"
  "path/filepath"

  "github.com/gorilla/mux"
  "github.com/bd878/live-connections/meta"
)

var publicPath = filepath.Join("../", "public")

type Server struct {
  manager *Manager
  httpServer *http.Server
  ctx context.Context
}

func NewHTTPServer(addr string, done chan struct{}) *Server {
  router := mux.NewRouter()
  manager := NewManager()

  router.HandleFunc("/ws/{area}/{user}", manager.HandleWS).Methods("GET")
  router.HandleFunc("/join", manager.HandleJoinArea).Methods("POST")
  router.HandleFunc("/area/new", manager.HandleNewArea).Methods("GET")
  router.HandleFunc("/area/{id}", manager.HandleAreaUsers).Methods("GET")

  httpServer := &http.Server{
    Addr: addr,
    Handler: router,
    IdleTimeout: 5 * time.Minute,
    ReadHeaderTimeout: time.Minute,
  }

  server := &Server{
    manager: manager,
    httpServer: httpServer,
    ctx: context.Background(),
  }

  go server.sigHandler(os.Interrupt, done)

  return server
}

func (s *Server) ListenAndServe() error {
  s.manager.StartHandlers(s.ctx)

  return s.httpServer.ListenAndServe()
}

func (s *Server) ListenAndServeTLS(serverCrt string, serverKey string) error {
  s.manager.StartHandlers(s.ctx)

  return s.httpServer.ListenAndServeTLS(serverCrt, serverKey)
}

func (s *Server) sigHandler(sig os.Signal, done chan struct{}) {
  sigint := make(chan os.Signal, 1)
  signal.Notify(sigint, sig)
  <-sigint

  if err := s.httpServer.Shutdown(context.Background()); err != nil {
    meta.Log().Debug(fmt.Sprintf("HTTP server Shutdown: %v", err))
  } else {
    meta.Log().Debug("SIGINT caught")
  }

  close(done)
}