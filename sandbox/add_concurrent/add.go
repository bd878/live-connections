package main

import (
  "math/rand"
  "sync/atomic"
  "sync"
)

func addSequential(numbers []int) int {
  var v int
  for _, n := range numbers {
    v += n
  }
  return v
}

func addConcurrent(goroutines int, numbers []int) int {
  var v int64
  totalNumbers := len(numbers)
  lastGoroutine := goroutines - 1
  stride := totalNumbers / goroutines

  var wg sync.WaitGroup
  wg.Add(goroutines)

  for g := 0; g < goroutines; g++ {
    go func(g int) {
      start := g * stride
      end := start + stride
      if g == lastGoroutine {
        end = totalNumbers
      }

      var lv int
      for _, n := range numbers[start:end] {
        lv += n
      }

      atomic.AddInt64(&v, int64(lv))
      wg.Done()
    }(g)
  }

  wg.Wait()

  return int(v)
}

func genNumbers(count int) []int {
  result := make([]int, count)
  for ; count > 0; count-- {
    result[count-1] = rand.Int()
  }
  return result
}
