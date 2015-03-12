# Pascal's triangle

GO implementation of [Pascal's triangle](http://en.wikipedia.org/wiki/Pascal%27s_triangle).


```
1 1 1 1 1 1
1 2 3 4 5 0
1 3 6 10 0 0
1 4 10 0 0 0
1 5 0 0 0 0
1 0 0 0 0 0
```

## Iterative solution

Generate all the values from top left to bottom right. Uses the previously generated values to compute the next.
This solution is optimal to generate the whole triangle.

```
for i := range output {
	for j := range output[i] {
		if i == 0 || j == 0 {
			output[i][j] = 1
			continue
		}
		output[i][j] = output[i-1][j] + output[i][j-1]
	}
}
```

## Recursive solution

Compute every single value independently. This solution is better to generate single values because it doesn't need to have other values in memory.

```
output[i][j] = innerRecursion(i, j)

func innerRecursion(i, j int) int {
	if i == 0 || j == 0 {
		return 1
	}
	return recursiveValue(i-1, j) + recursiveValue(i, j-1)
}
```