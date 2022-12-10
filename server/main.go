package main

import (
  "net/http"
  "io"
  "time"
  "log"
  "os"
  "fmt"
  "path/filepath"
  "os/signal"
  "context"
  "strings"

  "google.golang.org/grpc"

  pb "github.com/teralion/live-connections/server/proto"
)

const (
  serverAddr = "localhost:50051"
  diskRequestTimeout = 10*time.Second
)

var publicPath = filepath.Join("../", "public")

func main() {
  mux := http.NewServeMux()

  opts := []grpc.DialOption{grpc.WithInsecure()}

  conn, err := grpc.Dial(serverAddr, opts...)
  if err != nil {
    log.Fatalf("failed to dial: %v\n", err)
  }
  defer conn.Close()
  areaClient := pb.NewAreaManagerClient(conn)
  userClient := pb.NewUserManagerClient(conn)

  hub := GetHub()
  go hub.run()

  mux.Handle("/",
    http.HandlerFunc(
      func(w http.ResponseWriter, r *http.Request) {
        http.ServeFile(w, r, filepath.Join(publicPath, "index.html"))
      },
    ),
  )
  mux.Handle("/join",
    http.HandlerFunc(
      func(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodPost {
          http.Error(w,
            fmt.Sprintf("method %v not allowed\n", r.Method),
            http.StatusMethodNotAllowed,
          )
          return
        }
        ctx, cancel := context.WithTimeout(context.Background(), diskRequestTimeout)
        defer cancel()
        body, err := io.ReadAll(r.Body)
        if err != nil {
          http.Error(w,
            fmt.Sprint("cannot read body"),
            http.StatusBadRequest,
          )
        }
        var areaName strings.Builder
        areaName.Write(body)

        resp, err := userClient.Add(ctx, &pb.AddUserRequest{Area: areaName.String()})
        if err != nil {
          log.Fatalf("userClient.Add failed: %v", err)
        }
        fmt.Fprintf(w, "%v\n", resp.Name)
      },
    ),
  )
  mux.Handle("/area/new",
    http.HandlerFunc(
      func(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodGet {
          http.Error(w,
            fmt.Sprintf("method %v not allowed\n", r.Method),
            http.StatusMethodNotAllowed,
          )
          return
        }
        ctx, cancel := context.WithTimeout(context.Background(), diskRequestTimeout)
        defer cancel()
        resp, err := areaClient.Create(ctx, &pb.CreateAreaRequest{})
        if err != nil {
          log.Fatalf("areaClient.Create failed: %v", err)
        }
        fmt.Fprintln(w, resp.Name)
      },
    ),
  )
  mux.Handle("/area/",
    http.StripPrefix("/area/",
      http.HandlerFunc(
        func(w http.ResponseWriter, r *http.Request) {
          if r.Method != http.MethodGet {
            http.Error(w,
              fmt.Sprintf("method %v not allowed\n", r.Method),
              http.StatusMethodNotAllowed,
            )
            return
          }
          p := r.URL.Path;
          fmt.Printf("area path: %v\n", p)

          ctx, cancel := context.WithTimeout(context.Background(), diskRequestTimeout)
          defer cancel()
          resp, err := areaClient.ListUsers(ctx, &pb.ListAreaUsersRequest{Name: p})
          if err != nil {
            log.Fatalf("areaClient.ListUsers failed: %v", err)
          }
          fmt.Printf("% v\n", resp.GetUsers())
          fmt.Fprintf(w, "%v\n", resp.GetUsers())
        },
      ),
    ),
  )
  mux.Handle("/public/",
    http.StripPrefix("/public/",
      http.FileServer(http.Dir(publicPath)),
    ),
  )

  srv := &http.Server{
    Addr: "localhost:8080",
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
