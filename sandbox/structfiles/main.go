package main

import "fmt"

type Indexed struct {
  ID int
}

type Texted struct {
  Value string
}

type Test struct {
  *Texted
  *Indexed
}

func (t *Test) String() string {
  return fmt.Sprintf("%s - %d", t.Value, t.ID)
}

func main() {
  i := &Indexed{
    ID: 111,
  }

  v := &Texted{
    Value: "aaa",
  }

  t := &Test{
    v,
    i,
  }

  fmt.Println(t)
}