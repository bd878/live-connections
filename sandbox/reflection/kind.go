package main

import (
  "fmt"
  "reflect"
)

func ShowKinds() {
  for _, v := range []any{"Bangkok", 34, func(){}, []string{"a", "b"}} {
    switch v := reflect.ValueOf(v); v.Kind() {
    case reflect.String:
      fmt.Println(v.String())
    case reflect.Func:
      fmt.Println("func")
    case reflect.Int:
      fmt.Println(v.Int())
    default:
      fmt.Println("unhandleded kind %s", v.Kind())
    }
  }
}