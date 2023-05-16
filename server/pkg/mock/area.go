package mock

type Area struct {
  name string
  broadcast chan []byte
  clients map[string]Named
}

func NewArea() *Area {
  return &Area{
    name: "test",
    broadcast: make(chan []byte, 256),
    clients: make(map[string]Named, 1),
  }
}

func (p *Area) Name() string {
  return p.name
}

func (p *Area) SetName(n string) {
  p.name = name
}

func (p *Area) Broadcast() chan []byte {
  return p.broadcast
}

func (p *Area) Join(v interface{}) {
  n, ok := v.(Named)
  if !ok {
    panic("not named")
  }

  p.clients[n.Name()] = n
}

func (p *Area) Lose(v interface{}) {
  n, ok := v.(Named)
  if !ok {
    panic("not named")
  }

  delete(p.clients, n.Name())
}

func (p *Area) List() []string {
  var result []string
  for k, _ := range p.clients {
    result = append(result, k)
  }
  return result
}