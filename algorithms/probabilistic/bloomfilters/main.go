package main

import (
	"fmt"
	"os"
	"runtime/pprof"

	"github.com/willf/bloom"

	"code.google.com/p/go-uuid/uuid"
)

const MAX = 10000000

func main() {
	bf := bloom.New(20*MAX, 5) // load of 20, 5 keys

	hitCount := 0
	for i := 0; i < MAX; i++ {
		if (i+1)%1000 == 0 {
			fmt.Printf("Added %d elements, %d false hits.\n", i, hitCount)
		}
		if bf.TestAndAddString(uuid.New()) {
			hitCount++
		}
	}
	fmt.Printf("Miss rate with %d elements was %f. Estimated was %f.\n",
		MAX, (float64(hitCount) * 100.0 / MAX), bf.EstimateFalsePositiveRate(MAX))

	f, err := os.Create("bloom_mem.pprof")
	if err != nil {
		fmt.Println("Error: ", err)
	}
	pprof.WriteHeapProfile(f)
	f.Close()
}
