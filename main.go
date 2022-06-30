package main

import (
	// "github.com/antonyuhnovets/examples/concurrency"
	// "github.com/antonyuhnovets/examples/stability"
	"github.com/antonyuhnovets/examples/testchain"
	// "github.com/antonyuhnovets/examples/generics"
)

func main() {

	// generics.Exec()
	// concurrency.Exec_done()
	// concurrency.Exec_gen()

	// Using Circuit Breaker and Debounce patterns together:
	// wrapped := stability.Breaker(stability.Debounce(stability.Circuit))
	// responce, err = wrapped(ctx)

	// Using Retry pattern with emulator func:
	// wrapped := stability.Retry(stability.EmulateTransientError, 5, 2*time.Second)
	// result, err := wrapped(context.Background())

	// Using Timeout pattern:
	// ctx := context.Background()
	// ctxt, cancel := context.WithTimeout(ctx, 1*time.Second)
	// defer cancel()
	// timeout := stability.Timeout(stability.SlowFunc)
	// result, err := timeout(ctxt, "some input")

	// fmt.Println(result, err)

	testchain.Exec_chain()
}
