package main

import (
	"fmt"
	"time"
)

// Iterative solution
func PascalsTriangleIterative(size int) [][]int {
	output := makeMatrix(size)

	for i := range output {
		for j := range output[i] {
			if i == 0 || j == 0 {
				output[i][j] = 1
				continue
			}
			if i+j >= size {
				// Skip the second half of the matrix
				continue
			}
			output[i][j] = output[i-1][j] + output[i][j-1]
		}
	}
	return output
}

// Recursive solution
func PascalsTriangleRecursive(size int) [][]int {
	output := makeMatrix(size)

	for i := range output {
		for j := range output[i] {
			if i+j >= size {
				// Skip the second half of the matrix
				continue
			}
			output[i][j] = innerRecursion(i, j)
		}
	}
	return output
}

// Inner recursive function
func innerRecursion(i, j int) int {
	if i == 0 || j == 0 {
		return 1
	}
	return innerRecursion(i-1, j) + innerRecursion(i, j-1)
}

// Initialize the matrix
func makeMatrix(size int) [][]int {
	output := make([][]int, size)
	for i := range output {
		output[i] = make([]int, size)
	}
	return output
}

// Measures the duration
func bench(start time.Time) time.Duration {
	return time.Now().Sub(start)
}

func print(value [][]int) {
	for i := range value {
		fmt.Printf("%v\n", value[i])
	}
}

func main() {
	var output [][]int
	var t time.Time
	var d time.Duration

	t = time.Now()
	output = PascalsTriangleIterative(20)
	d = bench(t)
	fmt.Printf("Iterative solution %v\n", d)
	print(output)

	t = time.Now()
	output = PascalsTriangleRecursive(20)
	d = bench(t)
	fmt.Printf("Recursive solution %v\n", d)
	print(output)
}
