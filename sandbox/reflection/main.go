package main

import (
  "reflect"
  "log"
)

type Country struct {
  capital string
  name string
}

func main() {
  var residence Country
  residence = Country{capital: "Paris", name: "France"}
  residenceType := reflect.TypeOf(residence)

  if _, err := residenceType.FieldByName("capital"); err {
    log.Println("has field 'capital'")
  } else {
    log.Println("has not field 'capital'")
  }

  if _, err := residenceType.FieldByName("nation"); err {
    log.Println("has field 'nation")
  } else {
    log.Println("has not field 'nation'")
  }
}