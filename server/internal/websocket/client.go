package websocket

import (
  "time"
  "log"
  "encoding/binary"
  "bytes"
  "net/http"
  "strings"
  "errors"

  lc "github.com/teralion/live-connections/server/internal/conn"
  ws "github.com/gorilla/websocket"
)

const (
  writeWait = 10 * time.Second

  pongWait = 60 * time.Second

  pingPeriod = (pongWait * 9) / 10 // a bit less than pongWait
)

var enc = binary.LittleEndian

var lenWidth = 2

const (
  authMessageType int8 = 1

  mouseMoveMessageType int8 = 2

  listClientsOnlineMessageType int8 = 3

  authOkMessageType int8 = 4
)

var newline = []byte{'\n'}

var upgrader = ws.Upgrader{
  HandshakeTimeout: 10 * time.Second,
  ReadBufferSize: 512,
  WriteBufferSize: 512,
  CheckOrigin: checkClientOrigin,
}

type Coords struct {
  xPos float32

  yPos float32
}

type Client struct {
  conn *ws.Conn

  lc *lc.LiveConnections

  hub *Hub

  send chan []byte

  area string

  name string

  coords Coords
}

func checkClientOrigin(r *http.Request) bool {
  origin := r.Host
  return strings.Contains(origin, "localhost")
}

func NewClient(w http.ResponseWriter, r *http.Request, hub *Hub, liveConn *lc.LiveConnections) {
  conn, err := upgrader.Upgrade(w, r, nil)

  if err != nil {
    log.Println("failed to upgrade connection to WebSocket: ", err)
    return;
  }

  client := &Client{conn: conn, lc: liveConn, hub: hub, send: make(chan []byte, 256)}

  go client.readLoop()
  go client.writeLoop()
}

func (c *Client) isAuthenticated() bool {
  return c.area != "" && c.name != ""
}

func (c *Client) takeAuthMessage(mr *bytes.Reader) error {
  var err error

  var areaSize uint16
  if err := binary.Read(mr, enc, &areaSize); err != nil {
    return err
  }

  areaBytes := make([]byte, areaSize)
  if err = binary.Read(mr, enc, &areaBytes); err != nil {
    return err
  }

  var userSize uint16
  if err = binary.Read(mr, enc, &userSize); err != nil {
    return err
  }

  userBytes := make([]byte, userSize)
  if err = binary.Read(mr, enc, &userBytes); err != nil {
    return err
  }

  area := string(areaBytes)
  user := string(userBytes)

  userExists := c.lc.HasUser(area, user)

  if !userExists {
    return errors.New("no user")
  }

  c.area = area
  c.name = user

  return nil
}

func (c *Client) takeMouseMoveMessage(mr *bytes.Reader) error {
  var err error

  var xPos float32
  if err = binary.Read(mr, enc, &xPos); err != nil {
    log.Println("failed to read xPos =",  err)
    return err
  }

  var yPos float32
  if err = binary.Read(mr, enc, &yPos); err != nil {
    log.Println("failed to read YPos =",  err)
    return err
  }

  // rewrite each time new coords received
  c.coords = Coords{xPos: xPos, yPos: yPos}

  log.Println("take coords =", xPos, yPos)
  return nil
}

func (c *Client) readLoop() {
  c.conn.SetReadLimit(MaxPayloadSize)
  c.conn.SetReadDeadline(time.Now().Add(pongWait))
  c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
  for {
    _, r, err := c.conn.NextReader()
    if err != nil {
      log.Println("NextReader err =", err)
      break
    }

    var size uint16
    if err = binary.Read(r, enc, &size); err != nil {
      log.Println("binary.Read size err =", err)
      break
    }

    log.Println("size =", size)
    message := make([]byte, size)
    if err = binary.Read(r, enc, message); err != nil {
      log.Println("binary.Read message err =", err)
      break
    }
    log.Println("message =", message)

    if size != uint16(len(message)) {
      log.Println("size != message size = ", size, len(message))
      break
    }

    mr := bytes.NewReader(message)
    var messageType int8
    if err = binary.Read(mr, enc, &messageType); err != nil {
      log.Println("failed to read message type")
      break
    }

    if !c.isAuthenticated() && (messageType != authMessageType) {
      log.Println("message type is not auth message")
      break
    }

    switch messageType {
    case authMessageType:
      log.Println("auth message")
      if err = c.takeAuthMessage(mr); err != nil {
        log.Println("failed to parse auth message =", err)
        break
      }

      c.send <- doTextMessage("ok", authOkMessageType) // TODO: token
      c.hub.register <- c

      defer func() {
        c.hub.unregister <- c
      }()
    case mouseMoveMessageType:
      log.Println("mouse move message")
      if err = c.takeMouseMoveMessage(mr); err != nil {
        log.Println("failed to parse mouse move message =", err)
        break
      }

      var mouseMoveMessage []byte = c.doMouseMoveMessage()
      log.Println("mouse move response message =", mouseMoveMessage)
      c.send <- mouseMoveMessage // broadcast
    default:
      log.Println("unknown message type =", messageType)
      break
    }
  }
}

func (c *Client) writeLoop() {
  ticker := time.NewTicker(pingPeriod)

  defer func() {
    c.conn.Close()
    ticker.Stop()
  }()

  for {
    select {
    case message := <-c.send:
      if !c.isAuthenticated() {
        log.Println("not authenticated")
        continue
      }

      c.conn.SetWriteDeadline(time.Now().Add(writeWait))

      writer, err := c.conn.NextWriter(ws.BinaryMessage)
      if err != nil {
        log.Println("obtaining next writer err =", err)
        return
      }
      writer.Write(message)

      n := len(c.send)
      for i := 0; i < n; i++ {
        message = <-c.send
        writer.Write(message)
      }

      if err := writer.Close(); err != nil {
        log.Printf("writer close err: %v\n", err)
        return
      }
    case <-ticker.C:
      c.conn.SetWriteDeadline(time.Now().Add(writeWait))
      if err := c.conn.WriteMessage(ws.PingMessage, nil); err != nil {
        log.Printf("ping message writing failed, err: %v\n", err)
        return
      }
    }
  }
}

func (h *Hub) ListClientsOnline() []string {
  var names []string
  for client := range h.clients {
    names = append(names, client.name)
  }
  return names
}

func DoClientsOnlineMessage(users []string) []byte {
  var err error

  typeBytes := new(bytes.Buffer)
  if err = binary.Write(typeBytes, enc, listClientsOnlineMessageType); err != nil {
    log.Println("error writing message type =", err)
    return []byte{}
  }

  usersBytes := new(bytes.Buffer)
  for _, user := range users {
    var size uint16 = uint16(len(user))
    if err = binary.Write(usersBytes, enc, size); err != nil {
      log.Println("error writing user name size =", err)
      return []byte{}
    }

    if err = binary.Write(usersBytes, enc, []byte(user)); err != nil {
      log.Println("error writing user name size =", err)
      return []byte{}
    }
  }

  size := typeBytes.Len() + usersBytes.Len()
  sizeBytes := new(bytes.Buffer)
  if err = binary.Write(sizeBytes, enc, uint16(size)); err != nil {
    log.Println("error writing total size =", err)
    return []byte{}
  }

  return bytes.Join(
    [][]byte{
      sizeBytes.Bytes(),
      typeBytes.Bytes(),
      usersBytes.Bytes(),
    },
    []byte{},
  )
}

func doTextMessage(text string, messageType int8) []byte {
  var err error
  message := []byte(text)

  result := new(bytes.Buffer)
  var size uint16 = uint16(len(message))
  if err = binary.Write(result, enc, size); err != nil {
    log.Println("error writing text message size =", err)
    return []byte{}
  }

  if err = binary.Write(result, enc, messageType); err != nil {
    log.Println("error writing text message type =", err)
    return []byte{}
  }

  if err = binary.Write(result, enc, message); err != nil {
    log.Println("error writing text message =", err)
    return []byte{}
  }

  res := result.Bytes()
  return res
}

func (c *Client) doMouseMoveMessage() []byte {
  var err error

  // TODO: pack to struct
  typeBytes := new(bytes.Buffer)
  if err = binary.Write(typeBytes, enc, mouseMoveMessageType); err != nil {
    log.Println("error writing mouse move type =", err)
    return []byte{}
  }

  userBytes := new(bytes.Buffer)
  if err = binary.Write(userBytes, enc, []byte(c.name)); err != nil {
    log.Println("error writing user name bytes =", err)
    return []byte{}
  }

  userSizeBytes := new(bytes.Buffer)
  if err = binary.Write(userSizeBytes, enc, uint16(userBytes.Len())); err != nil {
    log.Println("error writing user name bytes size =", err)
    return []byte{}
  }

  coordsBytes := new(bytes.Buffer)
  if err = binary.Write(coordsBytes, enc, c.coords.xPos); err != nil {
    log.Println("error writing user x coord =", err)
    return []byte{}
  }

  if err = binary.Write(coordsBytes, enc, c.coords.yPos); err != nil {
    log.Println("error writing user y coord =", err)
    return []byte{}
  }

  size := typeBytes.Len() + userSizeBytes.Len() + userBytes.Len() + coordsBytes.Len()
  sizeBytes := new(bytes.Buffer)
  if err = binary.Write(sizeBytes, enc, uint16(size)); err != nil {
    log.Println("error writing total size =", err)
    return []byte{}
  }

  return bytes.Join(
    [][]byte{
      sizeBytes.Bytes(),
      typeBytes.Bytes(),
      userSizeBytes.Bytes(),
      userBytes.Bytes(),
      coordsBytes.Bytes(),
    },
    []byte{},
  )
}