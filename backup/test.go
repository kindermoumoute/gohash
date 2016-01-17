package main

import (
	"fmt"
	"runtime"
	"time"
)

var (
	numCPU int
	pi, h  float64
)

func main() {
	//numCPU = runtime.NumCPU()
	numCPU = 4
	runtime.GOMAXPROCS(numCPU)
	n := int64(10000000000)
	h = 1.0 / float64(n)
	pi = 0.0

	start := time.Now()
	c := make(chan float64, numCPU)
	for i := 0; i < numCPU; i++ {
		go chunk(int64(i)*n/int64(numCPU), (int64(i)+1)*n/int64(numCPU), c)
	}
	for i := 0; i < numCPU; i++ {
		pi += <-c
	}
	elapsed := time.Since(start)
	fmt.Println(elapsed.Seconds())
	fmt.Println(pi)
}

func f(a float64) float64 {
	return 4.0 / (1.0 + a*a)
}

func chunk(start, end int64, c chan float64) {
	var sum float64 = 0.0

	for i := start; i < end; i++ {
		x := h * (float64(i) + 0.5)
		sum += f(x)
	}
	c <- sum * h
}
