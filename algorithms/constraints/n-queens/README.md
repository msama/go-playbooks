#N-queens

http://en.wikipedia.org/wiki/Eight_queens_puzzle

Checks the validity of each queen using [exact cover constraints](http://en.wikipedia.org/wiki/Exact_cover).

```
func (s *Board) Check() bool {
	rows := make([]int, s.Size)
	cols := make([]int, s.Size)
	diags := make([]int, 2*s.Size-1)
	revDiags := make([]int, 2*s.Size-1)

	for _, q := range s.Queens {
		if q == nil {
			break
		}
		// One queen per rows
		rows[q.Row]++
		if rows[q.Row] > 1 {
			return false
		}

		// One queen per columns
		cols[q.Col]++
		if cols[q.Col] > 1 {
			return false
		}

		// One queen per diagonal
		diags[q.Row-q.Col+s.Size-1]++
		if diags[q.Row-q.Col+s.Size-1] > 1 {
			return false
		}

		// One queen per reverse diagonal
		revDiags[q.Row+q.Col]++
		if revDiags[q.Row+q.Col] > 1 {
			return false
		}
	}
	return true
}
```

# Exploring new solutions

Start placing one queen per column and move them down line by line until a solution is found.


```
-------
Solution: 91
-------
0 0 3 0 0 0 0 0 
0 0 0 0 5 0 0 0 
0 2 0 0 0 0 0 0 
0 0 0 0 0 0 0 8 
0 0 0 0 0 6 0 0 
0 0 0 4 0 0 0 0 
0 0 0 0 0 0 7 0 
1 0 0 0 0 0 0 0 
-------
Solution: 92
-------
0 0 3 0 0 0 0 0 
0 0 0 0 0 6 0 0 
0 0 0 4 0 0 0 0 
0 2 0 0 0 0 0 0 
0 0 0 0 0 0 0 8 
0 0 0 0 5 0 0 0 
0 0 0 0 0 0 7 0 
1 0 0 0 0 0 0 0 
-------
All solutions explored: 92
```