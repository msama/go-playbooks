package main

import (
	"fmt"
)

func TopDownMergeSort(input []int) {
	output := make([]int, len(input))
	TopDownSplitMerge(input, output)
}

func TopDownSplitMerge(input []int, output []int) {
	middle := len(input) / 2
	end := len(input)

	if end < 2 {
		return
	}

	TopDownSplitMerge(input[0:middle], output[0:middle])
	TopDownSplitMerge(input[middle:end], output[middle:end])
	TopDownMerge(input, output)
	copy(input, output)
}

func TopDownMerge(input []int, output []int) {
	begin := 0
	middle := len(input) / 2
	end := len(input)
	i0 := 0
	i1 := middle

	// While there are elements in the left or right
	for j := 0; j < end; j++ {
		// If left head exists and is <= existing right head.
		if i0 < middle && (i1 >= end || input[i0] <= input[i1]) {
			output[j] = input[i0]
			i0 = i0 + 1
		} else {
			output[j] = input[i1]
			i1 = i1 + 1
		}
	}
}

func main() {
	sortable := []int{1, 6, 5, 2, 4, 9, 0, 5, 3, 4, 7, 8}
	fmt.Printf("Input: %v\n", sortable)

	// Ascending
	TopDownMergeSort(sortable)
	fmt.Printf("Ascending: %v\n", sortable)
}
