package main

import "time"
import "math/rand"
import "log"

const (
  SIGSTOP int8 = 18

  SIGKILL int8 = 9
)

type Signal struct {
  pid int
  sig int8
}

type Queue struct {
  ps chan int
  kill chan *Signal
  quit chan bool
  proc map[int]bool
  halt map[int]bool
}

func NewQueue() *Queue {
  return &Queue{
    ps: make(chan int),
    kill: make(chan *Signal),
    quit: make(chan bool),
    proc: make(map[int]bool),
    halt: make(map[int]bool),
  }
}

func (q *Queue) Run() {
  var idle time.Duration
  for {
    idle = time.Duration(rand.Intn(1e3)) * time.Millisecond
    select {
    case pid := <-q.ps:
      q.proc[pid] = true
    case s := <-q.kill:
      switch s.sig {
      case SIGKILL:
        delete(q.proc, s.pid)
        delete(q.halt, s.pid)
      case SIGSTOP:
        delete(q.proc, s.pid)
        q.halt[s.pid] = true
      default:
        log.Println("unknown signal =", s.sig)
      }
    case <-time.After(idle):
      log.Println("idle =", idle)
    case <-q.quit:
      log.Println("quit")
      return
    }
  }
}

func main() {
  q := NewQueue()

  go q.Run()
  time.Sleep(3 * time.Second)
  q.quit <- true
}