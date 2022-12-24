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
  "github.com/teralion/live-connections/server/internal/conn"
  ws "github.com/teralion/live-connections/server/internal/websocket"
)

var publicPath = filepath.Join("../", "public")

func NewHTTPServer(addr string, done chan struct{}) *http.Server {
  liveConn := conn.NewLiveConnections()
  router := mux.NewRouter()

  hub := ws.NewHub()
  go hub.Run()

  handleWS := func(w http.ResponseWriter, r *http.Request) {
    ws.NewClient(w, r, hub, liveConn)
  }

  router.HandleFunc("/ws", handleWS).Methods("GET")
  router.HandleFunc("/join", liveConn.HandleJoin).Methods("POST")
  router.HandleFunc("/area/new", liveConn.HandleNewArea).Methods("GET")
  router.HandleFunc("/area/{id}", liveConn.HandleAreaUsers).Methods("GET")

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