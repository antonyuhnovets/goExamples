package stability

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/antonyuhnovets/examples/testchain"
)

func TestBreaker(t *testing.T) {
	ctx := context.Background()
	f := Breaker(CircuitMock, 5)

	for i := 1; i <= 15; i++ {
		result, err := f(ctx)
		if err != nil {
			fmt.Println(err)
		} else {
			t.Fatal(result)
		}
		time.Sleep(time.Second)
	}
}

func CircuitMock(ctx context.Context) (string, error) {
	chain := testchain.MakeChain()

	chain.TurnOffDown()

	err := chain.SendRequest(ctx, "GET")
	if err != nil {
		return "fail", err
	} else {
		return "success", nil
	}
}
