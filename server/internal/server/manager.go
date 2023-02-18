package server

import (
  "net/http"

  "github.com/gorilla/mux"
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

  conn, err := UpgradeConnection(w, r)
  if err != nil {
    return
  }

  if !m.disk.HasUser(area, user) {
    w.WriteHeader(http.StatusBadRequest)
    return
  }

  var hub *Hub
  if m.hubs[area] == nil {
    hub = NewHub()
    m.hubs[area] = hub

    go hub.Run()
  }

  client := NewClient(conn, hub, area, user)
  go client.ReadLoop()
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