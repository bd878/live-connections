package messages

import "encoding/binary"

const (
  auth int8 = 1
  mouseMove int8 = 2
  clientsOnline int8 = 3
  // 4, 5, 6 - free to use
  squareInit int8 = 7
  squareMove int8 = 8
  text int8 = 9
  recordsList int8 = 10
  addRecord int8 = 11
  selectRecord int8 = 12
)

var enc = binary.LittleEndian
