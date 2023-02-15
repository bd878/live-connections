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
  "github.com/teralion/live-connections/server/internal/rpc"
  "github.com/teralion/live-connections/server/internal/websocket"
)

var publicPath = filepath.Join("../", "public")

func NewHTTPServer(addr string, done chan struct{}) *http.Server {
  disk := rpc.NewDisk()
  router := mux.NewRouter()

  hub := websocket.NewHub()
  go hub.Run()

  handleWS := func(w http.ResponseWriter, r *http.Request) {
    websocket.NewClient(w, r, hub, disk)
  }

  router.HandleFunc("/ws", handleWS).Methods("GET")
  router.HandleFunc("/join", disk.HandleJoin).Methods("POST")
  router.HandleFunc("/area/new", disk.HandleNewArea).Methods("GET")
  router.HandleFunc("/area/{id}", disk.HandleAreaUsers).Methods("GET")

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