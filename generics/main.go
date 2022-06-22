package generics

import (
	"fmt"
)

type Number interface {
	int64 | float64
}

func SumIntFloat[K comparable, V Number](m map[K]V) V {
	var sum V
	for _, value := range m {
		sum += value
	}
	return sum
}

func Exec() {
	// Initialize a map for the integer values
	ints := map[string]int64{
		"first":  34,
		"second": 12,
	}

	// Initialize a map for the float values
	floats := map[string]float64{
		"first":  35.98,
		"second": 26.99,
	}

	fmt.Printf("Generic Sums: %v and %v\n",
		SumIntFloat(ints),   // Could be with types : SumIntsOrFloats[string, int64](ints)
		SumIntFloat(floats)) // SumIntsOrFloats[string, float64](floats))
}
