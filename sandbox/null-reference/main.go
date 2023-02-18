package main

import "log"

type City struct {
  name string
}

type Country struct {
  capital *City
  name string
}

func main() {
  payBas := Country{capital: nil, name: "Pay-Bas"}

  log.Println(payBas.name)
}