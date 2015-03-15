package main

import (
	"fmt"
)

// Creates a new queen
func NewQueen(row, col int) *Queen {
	return &Queen{
		Row: row,
		Col: col,
	}
}

// Defines the position of a queen in the board.
type Queen struct {
	Row int
	Col int
}

func NewBoard(size int) *Board {
	b := &Board{
		Size:   size,
		Queens: make([]*Queen, size),
	}
	return b
}

type Board struct {
	Size   int
	Queens []*Queen
}

// Verify that all the existing queens
// meets the requirements.
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

// Prints the board on stdout
func (s *Board) Print() {
	vals := make([][]int, s.Size)
	for r := 0; r < s.Size; r++ {
		vals[r] = make([]int, s.Size)
	}
	for i, q := range s.Queens {
		if q == nil {
			break
		}
		vals[q.Row][q.Col] = i + 1
	}
	fmt.Println("-------")
	for r := 0; r < s.Size; r++ {
		for c := 0; c < s.Size; c++ {
			fmt.Printf("%d ", vals[r][c])
		}
		fmt.Println("")
	}
	fmt.Println("-------")
}

func Solve(size int) {
	board := NewBoard(size)
	solutions := 0
	c := 0

	// Columns loop
	for c < size {
		queen := board.Queens[c]

		// Add a new queen to the board
		if queen == nil {
			// Create and add a new queen
			queen = NewQueen(0, c)
			board.Queens[c] = queen

			if board.Check() {
				// If the new added queen match the board
				// and if it was in the last colum then the board is complete
				if c == size-1 {
					solutions++
					fmt.Printf("Solution: %d\n", solutions)
					board.Print()
					//continue
				} else {
					// The new added queen matches
					// Keep adding
					c++
					// move to the next column
					continue
				}
			}
		}

		// Rows loop
		for true {
			// Move the existing queen down
			if queen.Row == size-1 {
				// Reached the end of the board
				if queen.Col == 0 {
					// All the solutions have been explored
					fmt.Printf("All solutions explored: %d\n", solutions)
					//board.Print()
					return
				} else {
					// Remove the current queen and move back to the previous column
					board.Queens[c] = nil
					queen = nil
					c--
					// Roll back and iterate
					break
				}
			} else {
				// Move down
				queen.Row++
				if board.Check() {
					if c == size-1 {
						solutions++
						fmt.Printf("Solution: %d\n", solutions)
						board.Print()
						// We just found a solution.
						// Let's keep exploring this column
					} else {
						// The queen fits its new position
						// Le'ts move on to the next column
						c++
						break
					}
				}
			}
		}
	}
}

func main() {
	Solve(8)
}
