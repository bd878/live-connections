package server

import (
  "net/http"

  "github.com/gorilla/mux"
  "github.com/teralion/live-connections/server/internal/meta"
  "github.com/teralion/live-connections/server/internal/rpc"
)

type Manager struct {
  disk *rpc.Disk
  hubs map[string]*Hub
}

func NewManager() *Manager {
  disk := rpc.NewDisk()

  return &Manager{
    disk,
    make(map[string]*Hub),
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

  meta.Log().Debug("connection upgraded")

  if !m.disk.HasUser(area, user) {
    w.WriteHeader(http.StatusBadRequest)
    meta.Log().Warn("area/user not exists")
    return
  } else {
    meta.Log().Debug("area/user has found")
  }

  var hub *Hub
  if m.hubs[area] == nil {
    meta.Log().Debug("no hub is running, creating")
    hub = NewHub(m.disk)
    m.hubs[area] = hub

    go hub.Run()
  } else {
    meta.Log().Debug("hub is running already, take it")
    hub = m.hubs[area]
  }

  client := NewClient(conn, hub, area, user)
  go client.ReadLoop()
  go client.LifecycleLoop()
  go client.WriteLoop()
}

func (m *Manager) HandleJoinArea(w http.ResponseWriter, r *http.Request) {
  m.disk.HandleJoin(w, r)
}

func (m *Manager) HandleNewArea(w http.ResponseWriter, r *http.Request) {
  m.disk.HandleNewArea(w, r)
}

func (m *Manager) HandleAreaUsers(w http.ResponseWriter, r *http.Request) {
  m.disk.HandleAreaUsers(w, r)
}