#Bubble sort

This is a GO implementation of the bubble sort algorithm. 

```
func BubbleSort(value sort.Interface) {
	var swap bool
	for i := 0; i < value.Len()-1; i++ {
		swap = false
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
```


The function takes a `sort.Interface` as input and sorts it in place. It also works well combined with the `sort.Reverse` function for descending order. 

```
// Ascending
BubbleSort(sort.IntSlice(sortable))

// Descending
BubbleSort(sort.Reverse(sort.IntSlice(sortable)))
```

