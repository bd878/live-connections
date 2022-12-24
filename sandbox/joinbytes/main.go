package main

import (
  "bytes"
  "encoding/binary"
  "log"
)

var enc = binary.LittleEndian

func getStr() string {
  return "hello"
}

func main() {
  var t int8 = 2
  var str string = getStr()
  var err error

  tBytes := new(bytes.Buffer)
  if err = binary.Write(tBytes, enc, t); err != nil {
    log.Println(err)
    return
  }

  strBytes := new(bytes.Buffer)
  if err = binary.Write(strBytes, enc, []byte(str)); err != nil {
    log.Println(err)
    return
  }

  result := bytes.Join(
    [][]byte{tBytes.Bytes(), strBytes.Bytes()},
    []byte{},
  )

  log.Println("result =", result)
}