package server

import (
  "net/http"
  "context"
  "io"
  "fmt"
  "strings"
  "os"
  "sync"

  "github.com/gorilla/mux"
  "github.com/bd878/live-connections/meta"
  "github.com/bd878/live-connections/server/pkg/rpc"
  "github.com/bd878/live-connections/server/pkg/protocol"
)

type Manager struct {
  disk *rpc.Disk
  areas map[string]*protocol.Area
  queue chan *protocol.Client
  maxConns int
  handlersCount int
  wg sync.WaitGroup
}

func NewManager() *Manager {
  diskAddr := os.Getenv("LC_DISK_ADDR")
  disk := rpc.NewDisk(diskAddr)

  return &Manager{
    disk: disk,
    areas: make(map[string]*protocol.Area),
    queue: make(chan *protocol.Client, 100),
    maxConns: 100,
    handlersCount: 10,
  }
}

func (m *Manager) HandleWS(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  area := vars["area"]
  user := vars["user"]

  meta.Log().Debug("area, user =", area, user)

  conn, err := UpgradeConnection(w, r)
  if err != nil {
    meta.Log().Debug("failed to upgrade connection =", err)
    return
  }

  if len(m.queue) > m.maxConns {
    meta.Log().Warn("queue is busy")
    return
  }

  u, err := m.disk.HasUser(context.Background(), area, user)
  if !u || err != nil {
    meta.Log().Warn("area/user not exists, break")
    return
  }

  cn := protocol.NewConn(conn)
  cn.Init()

  c := protocol.NewClient(cn)
  c.SetName(user)
  c.SetDisk(m.disk)

  var parent *protocol.Area
  if m.areas[area] == nil {
    parent = protocol.NewArea(area)
    parent.SetDisk(m.disk)

    m.areas[area] = parent

    go parent.Run(context.Background())
  } else {
    parent = m.areas[area]
  }

  c.SetParent(parent)

  m.queue <- c
}

func (m *Manager) Close() {
  close(m.queue)
}

func (m *Manager) StartHandlers(ctx context.Context) {
  meta.Log().Debug(fmt.Sprintf("launching %d goroutines", m.handlersCount))

  for i := 0; i < m.handlersCount; i++ {
    m.wg.Add(1)
    go m.Handle(ctx, i)
  }

  go func() {
    m.wg.Wait()
    m.Close()
  }()
}

func (m *Manager) Handle(ctx context.Context, num int) {
  go func() {
    <-ctx.Done()
    m.Close()
  }()

  for c := range m.queue {
    func() {
      meta.Log().Debug(fmt.Sprintf("handler %d get client %s", num, c.Name()))
      defer meta.Log().Debug(fmt.Sprintf("handler %d free client %s", num, c.Name()))

      c.Run(ctx)
    }()
  }
}

func (m *Manager) HandleJoinArea(w http.ResponseWriter, r *http.Request) {
  ctx := context.Background()

  body, err := io.ReadAll(r.Body)
  if err != nil {
    http.Error(w,
      fmt.Sprint("cannot read body"),
      http.StatusBadRequest,
    )
  }
  var areaName strings.Builder
  areaName.Write(body)

  userName, err := m.disk.CreateNewUser(ctx, areaName.String())
  if err != nil {
    meta.Log().Error("failed to create user", err)
    http.Error(w,
      fmt.Sprint("cannot create user"),
      http.StatusBadRequest,
    )
    return
  }

  record, err := m.disk.AddTextRecord(ctx, areaName.String(), userName)
  if err != nil {
    meta.Log().Error("failed to add text record to new user", err)
    http.Error(w,
      fmt.Sprint("cannot create user"),
      http.StatusInternalServerError,
    )
    return
  }
  err = m.disk.SelectTextRecord(ctx, areaName.String(), userName, record.ID)
  if err != nil {
    meta.Log().Error("failed to select new text record for new user", err)
    http.Error(w,
      fmt.Sprint("cannot create user"),
      http.StatusInternalServerError,
    )
    return
  }

  w.Header().Set("Content-Type", "text/plain; charset=utf-8")
  w.Header().Set("Access-Control-Allow-Origin", "*")
  fmt.Fprint(w, userName)
}

func (m *Manager) HandleNewArea(w http.ResponseWriter, r *http.Request) {
  ctx := context.Background()

  areaName, err := m.disk.CreateNewArea(ctx)
  if err != nil {
    meta.Log().Error(fmt.Sprintf("area.Create failed: %v", err))
    http.Error(w,
      fmt.Sprint("cannot create area"),
      http.StatusBadRequest,
    )
    return
  }

  w.Header().Set("Content-Type", "text/plain; charset=utf-8")
  w.Header().Set("Access-Control-Allow-Origin", "*")
  fmt.Fprint(w, areaName)
}

func (m *Manager) HandleAreaUsers(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  area := vars["id"]

  ctx := context.Background()

  users, err := m.disk.ListUsers(ctx, area)
  if err != nil {
    meta.Log().Error(fmt.Sprintf("area.ListUsers failed: %v", err))
    http.Error(w,
      fmt.Sprint("cannot create area"),
      http.StatusBadRequest,
    )
    return
  }

  w.Header().Set("Content-Type", "text/plain; charset=utf-8")
  w.Header().Set("Access-Control-Allow-Origin", "*")
  fmt.Fprint(w, users)
}