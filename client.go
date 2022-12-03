package main

import (
	"time"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

const (
	writeWait = 10 * time.Second

	pongWait = 60 * time.Second

	maxRead = 512 // bytes

	pingPeriod = (pongWait * 9) / 10 // a bit less than pongWait
)

var (
	newline = []byte{'\n'}
)

var upgrader = websocket.Upgrader{
	HandshakeTimeout: 10 * time.Second,
	ReadBufferSize: 1024,
	WriteBufferSize: 1024,
}

type Client struct {
	conn *websocket.Conn

	hub *Hub

	send chan []byte
}

func NewClient(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Println("failed to upgrade connection to WebSocket: ", err)
		return;
	}

	hub := GetHub()
	if err != nil {
		log.Printf("error obtaining hub: %v\n", err)
		return
	}

	client := &Client{conn: conn, hub: hub, send: make(chan []byte, 256)}
	hub.register <- client

	// go client.readLoop()
	go client.writeLoop()
}

func (c *Client) writeLoop() {
	ticker := time.NewTicker(pingPeriod)

	defer func() {
		c.conn.Close()
		ticker.Stop()
		c.hub.unregister <- c
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			writer, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				log.Printf("obtaining next writer err: %v\n", err)
				return
			}
			writer.Write(message)

			n := len(c.send)
			for i := 0; i < n; i++ {
				writer.Write(newline)
				writer.Write(<-c.send)
			}

			if err := writer.Close(); err != nil {
				log.Printf("writer close err: %v\n", err)
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Printf("ping message writing failed, err: %v\n", err)
				return
			}
		}
	}
}