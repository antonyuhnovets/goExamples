package main

import (
	"github.com/antonyuhnovets/examples/concurrency"
	// "github.com/antonyuhnovets/examples/stability"
	// "github.com/antonyuhnovets/examples/generics"
)

func main() {

	// generics.Exec()
	// concurrency.Exec_done()

	// Using circuit breaker and debounce patterns together
	// wrapped = stability.Breaker(stability.Debounce(stability.Circuit))
	// responce, err = wrapped(ctx)

	concurrency.Exec_gen()
}
