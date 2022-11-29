package main

import (
  "net/http"
  "io"
  "time"
  "log"
  "os"
  "os/signal"
  "context"
)

func main() {
  srv := &http.Server{
    Addr: ":8080",
    Handler: EchoHandler(),
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

  if err := srv.ListenAndServe(); err != http.ErrServerClosed {
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
