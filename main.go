package main

import (
  "net/http"
  "io"
  "time"
  "log"
  "os"
  "path/filepath"
  "os/signal"
  "context"

  "github.com/gorilla/websocket"
)

func main() {
  mux := http.NewServeMux()

  var upgrader = websocket.Upgrader{
    ReadBufferSize: 1024,
    WriteBufferSize: 1024,
  }

  mux.Handle("/",
    http.HandlerFunc(
      func(w http.ResponseWriter, r *http.Request) {
        http.ServeFile(w, r, filepath.Join("public/", "index.html"))
      },
    ),
  )
  mux.Handle("/public/",
    http.StripPrefix("/public/",
      http.FileServer(http.Dir("public/")),
    ),
  )
  mux.Handle("/ws",
    http.HandlerFunc(
      func(w http.ResponseWriter, r *http.Request) {
        conn, err := upgrader.Upgrade(w, r, nil)
        if err != nil {
          log.Println(err)
          return
        }

        for {
          messageType, p, err := conn.ReadMessage()
          if err != nil {
            log.Println(err)
            return
          }
          if err := conn.WriteMessage(messageType, p); err != nil {
            log.Println(err)
            return
          }
        }
      },
    ),
  )

  srv := &http.Server{
    Addr: ":8080",
    Handler: mux,
    IdleTimeout: 5 * time.Minute,
    ReadHeaderTimeout: time.Minute,
  }

  done := make(chan struct{})
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

  log.Println("Server is listening on :8080")
  if err := srv.ListenAndServeTLS("server.crt", "server.key"); err != http.ErrServerClosed {
    log.Fatalf("HTTP server ListenAndServe: %v", err)
  }

  <-done
}

func EchoHandler() http.Handler {
  return http.HandlerFunc(
    func (w http.ResponseWriter, r *http.Request) {
      b, err := io.ReadAll(r.Body)
      if err != nil {
        http.Error(w, "Internal server error",
          http.StatusInternalServerError)
        return
      }

      _, _ = w.Write(append(b, '\n'))
    })
}
