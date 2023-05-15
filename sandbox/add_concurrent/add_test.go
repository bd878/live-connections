package main

import (
  "runtime"
  "testing"
)

var numbers = genNumbers(10e6)

func BenchmarkSequential(b *testing.B) {
  for i := 0; i < b.N; i++ {
    addSequential(numbers)
  }
}

func BenchmarkConcurrent(b *testing.B) {
  for i := 0; i < b.N; i++ {
    addConcurrent(runtime.NumCPU(), numbers)
  }
}