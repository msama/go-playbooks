package main

import (
	"fmt"
	"math/rand"
	"time"
)

func NewGeneratorCMA(r *rand.Rand) *GeneratorCMA {
	return &GeneratorCMA{
		Source:  r,
		Average: 0,
		Count:   0,
	}
}

// Generator which keeps a total counter used to approximate the average
// Alas if the stream is really infinite the counter will eventually overflow
type GeneratorCMA struct {
	Source  *rand.Rand
	Average float64
	Count   float64
}

func (s *GeneratorCMA) Next() float64 {
	val := s.Source.ExpFloat64()
	s.Average = s.Average*(s.Count/(s.Count+1)) + val/(s.Count+1)
	s.Count++
	return val
}

func main() {
	var total float64 = 0
	count := 10000000

	generator := NewGeneratorCMA(rand.New(rand.NewSource(time.Now().UnixNano())))

	for i := 0; i < count; i++ {
		total += generator.Next()
	}

	totalAverage := total / float64(count)
	fmt.Printf("%e <- Average\n", totalAverage)
	fmt.Printf("%e <- Moving average (total). Error: %e\n",
		generator.Average, generator.Average-totalAverage)
}
