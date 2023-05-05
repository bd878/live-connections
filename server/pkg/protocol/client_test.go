package protocol

import (
  "testing"
  "io"
  "fmt"
  "time"
  "log"
  "context"
  "sync"
  "os"
  "os/signal"
  "net/http"
)

type NullConn struct {}

func (c *NullConn) NextReader() (int, io.Reader, error) {
  return 0, nil, nil
}

func (c *NullConn) WriteMessage(_ int, _ []byte) error {
  return nil
}

func (c *NullConn) Close() error {
  return nil
}

func TestClient(_ *testing.T) {
  var wg sync.WaitGroup

  f := func(w http.ResponseWriter, r *http.Request) {
    fmt.Println("got connection")

    client := NewClient(&NullConn{})
    client.SetName("test")
    fmt.Println(client.Name())
  }

  serv := &http.Server{
    Addr: "localhost:8080",
    Handler: http.HandlerFunc(f),
  }

  go func() {
    c := make(chan os.Signal, 1)
    signal.Notify(c, os.Interrupt)

    <-c
    log.Println("got signal")
    if err := serv.Shutdown(context.Background()); err != http.ErrServerClosed {
      log.Fatal(err)
    }
  }()

  wg.Add(1)
  go func() {
    defer wg.Done()
    fmt.Println("server is listening...") 
    if err := serv.ListenAndServe(); err != nil {
      log.Fatal(err)
    }
  }()

  t := time.NewTimer(1500*time.Millisecond)
  <-t.C

  _, err := http.Get("http://localhost:8080")
  if err != nil {
    log.Fatal(err)
  }

  wg.Wait()
}
