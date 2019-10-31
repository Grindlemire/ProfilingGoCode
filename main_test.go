package main

import (
	"fmt"
	"image"
	"log"
	"testing"
)

var tests = []struct {
	Name string
	f    func(opts Opts) (img image.Image, err error)
}{
	{"basic", executeAlgorithm},
	{"pixelParallel", executePixelParallelAlgorithm},
	{"columnParallel", executeColumnParallelAlgorithm},
	{"workers", executeWorkersAlgorithm},
	{"bufferedWorkers", executeBufferedWorkersAlgorithm},
	{"bufferedColumnsWorkers", executeBufferedColumnWorkersAlgorithm},
}

func BenchmarkAlgorithms(b *testing.B) {
	for n := 0; n <= 500; n += 50 {
		opts := Opts{
			Complexity: 40,
			Width:      n,
			Height:     n,
			MoveX:      -0.75,
		}
		for _, test := range tests {
			b.Run(fmt.Sprintf("%s/%d", test.Name, n), func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					_, err := test.f(opts)
					if err != nil {
						log.Fatalf("Error: %s", err)
					}
				}
			})
		}
	}
}
