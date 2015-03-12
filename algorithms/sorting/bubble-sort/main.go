package main

import (
	"fmt"
	"sort"
)

func BubbleSort(value sort.Interface) {
	// If no element are swapped in the inner loop
	// then the collection is already sorted.
	var swap bool
	// Main loop
	for i := 0; i < value.Len()-1; i++ {
		swap = false
		// Inner loop
		for j := i; j < value.Len(); j++ {
			if value.Less(j, i) {
				value.Swap(i, j)
				swap = true
			}
		}
		// No swap was done in the inner circle
		if !swap {
			return
		}
	}
}

func main() {
	sortable := []int{1, 6, 5, 2, 4, 9, 0, 5, 3, 4, 7, 8}
	fmt.Printf("Input: %v\n", sortable)

	// Ascending
	BubbleSort(sort.IntSlice(sortable))
	fmt.Printf("Ascending: %v\n", sortable)

	// Descending
	BubbleSort(sort.Reverse(sort.IntSlice(sortable)))
	fmt.Printf("Descending: %v\n", sortable)
}
