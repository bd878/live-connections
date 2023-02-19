package main

import "fmt"

type A uint64

func (a *A) Set(value A) {
  *a = value
}

func main() {
  var a A
  a.Set(A(15))
  fmt.Println(a)
}