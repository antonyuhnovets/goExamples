package testchain

import (
	"context"
	// 	"github.com/antonyuhnovets/examples/stability"
)

func MakeGetRequest(ctx context.Context) map[Request]Response {
	chain := MakeChain()

	chain.SendRequest(ctx, "GET")

	return chain.up.srv.storage
}

func CircuitT(ctx context.Context) (string, error) {
	chain := MakeChain()

	chain.TurnOffDown()

	err := chain.SendRequest(ctx, "GET")

	if err != nil {
		return "fail", err
	} else {
		return "success", nil
	}
}

func Exec_chain() {
	return
	// ctx := context.Background()

	// f := stability.Breaker(CircuitT, 5)

	// f(ctx)
	// s := MakeGetRequest(ctx)

	// for k, v := range s {
	// 	fmt.Printf("Request : %+v; Response: %+v", k, v)
	// }
}
