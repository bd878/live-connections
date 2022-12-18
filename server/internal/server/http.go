package server

import (
  "net/http"
  "time"
  "os"
  "context"
  "os/signal"
  "log"
  "path/filepath"

  "github.com/gorilla/mux"
)

var publicPath = filepath.Join("../", "public")

func NewHTTPServer(addr string, done chan struct{}) *http.Server {
  liveConn := NewLiveConnections()
  router := mux.NewRouter()

  router.HandleFunc("/ws", liveConn.handleWS).Methods("GET")
  router.HandleFunc("/join", liveConn.handleJoin).Methods("POST")
  router.HandleFunc("/area/new", liveConn.handleNewArea).Methods("GET")
  router.HandleFunc("/area/{id}", liveConn.handleAreaUsers).Methods("GET")

  srv := &http.Server{
    Addr: addr,
    Handler: router,
    IdleTimeout: 5 * time.Minute,
    ReadHeaderTimeout: time.Minute,
  }

  go func() {
    sigint := make(chan os.Signal, 1)
    signal.Notify(sigint, os.Interrupt)
    <-sigint

    if err := srv.Shutdown(context.Background()); err != nil {
      log.Printf("HTTP server Shutdown: %v", err)
    } else {
      log.Println("SIGINT caught")
    }
    close(done)
  }()

  return srv
}